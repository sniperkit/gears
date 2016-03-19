package gears

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"golang.org/x/net/context"
)

// BGContext is the background context for all
// gears middleware
var BGContext context.Context

func init() {
	BGContext = context.Background()
}

// ContextHandler is a function signature for handers which require context
type ContextHandler func(c context.Context, w http.ResponseWriter, r *http.Request)

// Gear is a context aware middleware function signature
type Gear func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context

// New Gear is constructed by taking either of the following types as input:
//
// - func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context
//
// - func(c context.Context, w http.ResponseWriter, r *http.Request)
//
// - http.Handler
//
// - http.HandlerFunc
//
// Passing other types will panic.
func New(fn interface{}) Gear {
	switch t := fn.(type) {
	case func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context:
		return Gear(t)
	case func(c context.Context, w http.ResponseWriter, r *http.Request):
		return wrapContextHandler(t)
	case http.Handler:
		return wrapHandler(t)
	case http.HandlerFunc:
		return wrapHandlerFunc(t)
	default:
		panic(fmt.Sprintf("invalid parameter type for gears.New (%v)\n", reflect.TypeOf(fn).Kind()))
	}
}

func wrapHandler(h http.Handler) Gear {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		h.ServeHTTP(w, r)
		return c
	}
}

func wrapContextHandler(h ContextHandler) Gear {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		h(c, w, r)
		return c
	}
}

func wrapHandlerFunc(fn http.HandlerFunc) Gear {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		fn(w, r)
		return c
	}
}

// loggedWriter
type loggedWriter struct {
	status int
	w      http.ResponseWriter
}

func (lw *loggedWriter) WriteHeader(status int) {
	lw.status = status
	lw.w.WriteHeader(status)
}

func (lw *loggedWriter) Header() http.Header {
	return lw.w.Header()
}

func (lw *loggedWriter) Write(b []byte) (int, error) {
	return lw.w.Write(b)
}

// Chain multiple middleware returning a single Gear func.
func Chain(gears ...Gear) Gear {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		var localCtx context.Context
		for _, gear := range gears {
			localCtx = gear(c, w, r)
			if localCtx == nil {
				return NewErrorContext(c, NewStatusError(500, "Middleware returned nil context"))
			}
			if localCtx.Err() != nil {
				return localCtx
			}

			c = localCtx
		}

		return c
	}
}

func handleError(c context.Context, w http.ResponseWriter) {

	// handle http error
	errValue := c.Value("error")
	if errValue == nil {
		// error not found, means that a middleware canceled the context
		// signaling that the request is completed, don't do further processing
		return
	}

	res := detailedError{}
	switch t := errValue.(type) {
	case DetailedError:
		res.status = t.Status()
		res.ErrorCode = t.Code()
		res.ErrorDescription = t.Description()
		res.ErrorDetails = t.Details()
	case Error:
		res.status = t.Status()
		res.ErrorCode = t.Code()
		res.ErrorDescription = t.Description()
	case StatusError:
		res.status = t.Status()
		res.ErrorDescription = t.Error()
	default:
		res.status = 500
		res.ErrorCode = "unknown_error"
		res.ErrorDescription = "application returned invalid error message"
		res.ErrorDetails = errValue
	}

	responseBody, _ := json.Marshal(res)

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(res.Status())
	w.Write(responseBody)
}

// NewStatusError sets the error on the context and returns the canceled context.
// It is going to be deprecated, please use NewError instead which can create more
// precise error messages.
func NewStatusError(status int, description string) StatusError {
	return detailedError{status: status, ErrorDescription: description}
}

// NewErrorContext expects a context and err, the latter implementing StatusError interface.
// It returns a canceled context with the error set under "error" key.
func NewErrorContext(c context.Context, err StatusError) context.Context {

	var cancel context.CancelFunc
	c, cancel = context.WithCancel(c)
	defer cancel()

	return context.WithValue(c, "error", err)
}

// NewCanceledContext return a context which is canceled. It is used for signaling
// to a subsequent handler / gear / middleware in the chain to stop processing the request.
func NewCanceledContext(c context.Context) context.Context {

	var cancel context.CancelFunc
	c, cancel = context.WithCancel(c)
	cancel()
	return c

}

func (gear Gear) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c, cancel := context.WithCancel(BGContext)

	defer func() {
		cancel()
	}()

	c = gear(c, w, r)
	switch c.Err() {
	case context.Canceled, context.DeadlineExceeded:
		handleError(c, w)
		return
	}
}
