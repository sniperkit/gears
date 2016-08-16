# gears
--
    import "github.com/zgiber/gears"

## Important

With the introduction of context in go 1.7 the benefits of using this package are probably not worth
introducing external dependencies. Use the standard package's context, and the helper methods on the
http.Request instead.

## Usage

```go
var BGContext context.Context
```
BGContext is the background context for all gears middleware

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
NewErrorContext expects a context and err, the latter implementing StatusError
interface. It returns a canceled context with the error set under "error" key.

#### type ContextHandler

```go
type ContextHandler func(c context.Context, w http.ResponseWriter, r *http.Request)
```

ContextHandler is a function signature for handers which require context

#### type DetailedError

```go
type DetailedError interface {
	Error
	Details() interface{}
}
```

DetailedError is an interface for being used for cases where more context is
required to be provided with the error message.

#### type Error

```go
type Error interface {
	StatusError
	Code() string
	Description() string
}
```

Error is an interface used in cases when an error-code string and a longer
description needs to be passed with the original error message

#### func  NewError

```go
func NewError(status int, code, description string, details interface{}) Error
```
NewError returns an implementation of Error interface. Use NewError to set the
value for the "error" key in the context of a middleware. The details can be a
nil interface.

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
New Gear is constructed by taking either of the following types as input:

- func(c context.Context, w http.ResponseWriter, r *http.Request)
context.Context

- func(c context.Context, w http.ResponseWriter, r *http.Request)

- http.Handler

- http.HandlerFunc

Passing other types will panic.

#### func (Gear) ServeHTTP

```go
func (gear Gear) ServeHTTP(w http.ResponseWriter, r *http.Request)
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

StatusError is an interface used in case an http status code needs to be passed
with the error message

#### func  NewStatusError

```go
func NewStatusError(status int, description string) StatusError
```
NewStatusError sets the error on the context and returns the canceled context.
It is going to be deprecated, please use NewError instead which can create more
precise error messages.
