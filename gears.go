package gears

import (
	"encoding/json"
	"fmt"
	"net/http"

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

// Handler is a context aware http request handler
type Handler struct {
	fn   func(c context.Context, w http.ResponseWriter, r *http.Request)
	gear Gear
}

// Gear is a context aware middleware function signature
type Gear func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context

// httpError contains status code and message
// and implements error interface
type httpError struct {
	status  int
	message string
}

func (err *httpError) Error() string {
	return fmt.Sprintf("%v %s", err.Status(), err.message)
}

func (err *httpError) Status() int {
	return err.status
}

func (err *httpError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{"error": err.Status(), "description": err.Error()})
}

// NewHTTPError returns a httpError as an error interface
func newHTTPError(status int, message string) *httpError {
	return &httpError{status, message}
}

// NewHandler returns a pointer to a Handler struct which implements
// http.Handler interface. This is a convenient way to construct context aware
// gear.Handlers which can be used with standard http routers.
// fn must have a signature of either func(w http.ResponseWriter, r *http.Request)
// or func(c context.Context, w http.ResponseWriter, r *http.Request)
func NewHandler(fn interface{}, gears ...Gear) *Handler {
	var handlerFn ContextHandler
	switch t := fn.(type) {
	case func(c context.Context, w http.ResponseWriter, r *http.Request):
		handlerFn = t
	case func(w http.ResponseWriter, r *http.Request):
		handlerFn = withContext(t)
	case http.Handler:
		handlerFn = withContext(t.ServeHTTP)
	default:
		panic("invalid handler signature")
	}
	gear := Chain(gears...)
	return &Handler{handlerFn, gear}
}

// allows for using simple handlers (those without context in NewHandler)
func withContext(fn func(w http.ResponseWriter, r *http.Request)) ContextHandler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, cancel := context.WithCancel(BGContext)
	defer cancel()
	c = h.gear(c, w, r)
	switch c.Err() {
	case context.Canceled, context.DeadlineExceeded:
		handleError(c, w)
		return
	}
	h.fn(c, w, r)
}

// Chain multiple middleware
func Chain(gears ...Gear) Gear {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		var localCtx context.Context
		for _, gear := range gears {
			localCtx = gear(c, w, r)
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

	statusErr, ok := errValue.(StatusError)
	if !ok {
		// error doesn't implement StatusError
		statusErr = NewStatusError(500, fmt.Sprint(errValue))
	}

	responseBody, err := json.Marshal(statusErr)
	if err != nil {
		// error can't be marshaled
		statusErr = NewStatusError(500, fmt.Sprint(errValue))
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusErr.Status())
	fmt.Fprintln(w, string(responseBody))
}

// NewStatusError sets the error on the context and returns the canceled context.
func NewStatusError(status int, message string) StatusError {
	return &httpError{status, message}
}

// NewErrorContext expects an err which implements StatusError interface, and returns
// a context which has a json formatted error on it.
func NewErrorContext(c context.Context, err StatusError) context.Context {

	var cancel context.CancelFunc
	c, cancel = context.WithCancel(c)
	defer cancel()

	if jsonErr, ok := err.(JSONError); ok {
		return context.WithValue(c, "error", jsonErr)
	}

	return context.WithValue(c, "error", &httpError{err.Status(), err.Error()})
}

// NewCanceledContext return a context which is canceled. It is used for signaling
// to any subsequent handler / gear / middleware in the chain to stop processing the request.
func NewCanceledContext(c context.Context, err StatusError) context.Context {

	var cancel context.CancelFunc
	c, cancel = context.WithCancel(c)
	defer cancel()
	return c

}
