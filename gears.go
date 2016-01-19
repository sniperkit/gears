package gears

import (
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

// Handler is a context aware http request handler
type Handler struct {
	fn   func(c context.Context, w http.ResponseWriter, r *http.Request)
	gear Gear
}

// Gear is a context aware middleware function signature
type Gear func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context

// httpError contains code and message
// and implements error interface
type httpError struct {
	code    int
	message string
}

func (err *httpError) Error() string {
	return fmt.Sprintf("%v %s", err.code, err.message)
}

// NewHTTPError returns a httpError as an error interface
func newHTTPError(code int, message string) *httpError {
	return &httpError{code, message}
}

// NewHandler returns a pointer to a Handler struct which implements
// http.Handler interface. This is a convenient way to construct context aware
// gear.Handlers which can be used with standard http routers.
func NewHandler(fn func(c context.Context, w http.ResponseWriter, r *http.Request), gears ...Gear) *Handler {
	gear := Chain(gears...)
	return &Handler{fn, gear}
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
	if err, ok := c.Value("error").(*httpError); ok {
		http.Error(w, err.Error(), err.code)
	} else {
		http.Error(w, "wrong middleware error type", http.StatusInternalServerError)
	}
}

// NewError sets the error on the context and returns the canceled context.
func NewError(c context.Context, code int, message string) context.Context {

	var cancel context.CancelFunc // cancel the context
	err := &httpError{code, message}
	c = context.WithValue(c, "error", err)
	c, cancel = context.WithCancel(c)
	cancel()

	return c
}
