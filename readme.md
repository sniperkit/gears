# gears
--
    import "github.com/zgiber/gears"


## Usage

```go
var BGContext context.Context
```
BGContext is the background context for all gears middleware

#### func  NewErrorContext

```go
func NewErrorContext(c context.Context, err StatusError) context.Context
```
NewErrorContext expects an err which implements StatusError interface, and
returns a context which has a json formatted error on it.

#### type Gear

```go
type Gear func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context
```

Gear is a context aware middleware function signature

#### func  Chain

```go
func Chain(gears ...Gear) Gear
```
Chain multiple middleware

#### type Handler

```go
type Handler struct {
}
```

Handler is a context aware http request handler

#### func  NewHandler

```go
func NewHandler(fn func(c context.Context, w http.ResponseWriter, r *http.Request), gears ...Gear) *Handler
```
NewHandler returns a pointer to a Handler struct which implements http.Handler
interface. This is a convenient way to construct context aware gear.Handlers
which can be used with standard http routers.

#### func (*Handler) ServeHTTP

```go
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

#### type JSONError

```go
type JSONError interface {
	Error() string
	MarshalJSON() ([]byte, error)
}
```

JSONError is an interface used for handling errors as JSON encoded bytes.

#### type StatusError

```go
type StatusError interface {
	Error() string
	Status() int
}
```

StatusError is an interface used for handling http errors with proper statuses

#### func  NewStatusError

```go
func NewStatusError(status int, message string) StatusError
```
NewStatusError sets the error on the context and returns the canceled context.
