# gears
--
    import "github.com/zgiber/gears"


## Usage

```go
var BGContext context.Context
```
BGContext is the background context for all gears middleware. Usually the http.HandleFunc
Is provided by NewHandler which cancels the derived context as soon the http request is
handled.  

#### func  NewCanceledContext

```go
func NewCanceledContext(c context.Context) context.Context
```
NewCanceledContext return a context which is canceled. It is used for signaling
to a subsequent handler / gear / middleware in the chain to stop processing the
request.

#### func  NewErrorContext

```go
func NewErrorContext(c context.Context, err StatusError) context.Context
```
NewErrorContext expects an err which implements StatusError interface, and
returns a context which has a json formatted error on it.

#### func  NewHandler

```go
func NewHandler(logger Logger, gears ...Gear) http.Handler
```
NewHandler returns a http.Handler as a convenient way to construct context aware
gear.Handlers which can be used with standard http routers. fn must have a
signature of either func(w http.ResponseWriter, r *http.Request) or func(c
context.Context, w http.ResponseWriter, r *http.Request)

#### type ContextHandler

```go
type ContextHandler func(c context.Context, w http.ResponseWriter, r *http.Request)
```

ContextHandler is a function signature for handers which require context

#### type Gear

```go
type Gear func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context
```

Gear is a context aware middleware function signature

#### func  Chain

```go
func Chain(gears ...Gear) Gear
```
Chain multiple middleware returning a single Gear func.

#### func  New

```go
func New(fn interface{}) Gear
```
New Gear is constructed by taking either of the following types as input;  
func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context  
func(c context.Context, w http.ResponseWriter, r *http.Request)  
http.Handler  
http.HandlerFunc  
Passing other types will return error.

#### func  WrapHandlerFunc

```go
func WrapHandlerFunc(fn http.HandlerFunc) Gear
```
WrapHandlerFunc wraps a http.HandlerFunc returning a Gear. This provides a way
to combine common middleware packages with gears.

#### type JSONError

```go
type JSONError interface {
	Error() string
	MarshalJSON() ([]byte, error)
}
```

JSONError is an interface used for handling errors as JSON encoded bytes.

#### type Logger

```go
type Logger interface {
	Printf(format string, v ...interface{})
}
```

Logger is an interface which is used by gears to log return code and completion
time on each http request

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
