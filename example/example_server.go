package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zgiber/gears"
	"golang.org/x/net/context"
)

func middlewareExample(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	w.Write([]byte("Hello\n")) // do something (normally the middleware doesn't write to the response body, this is just for testing)
	return c
}

func errorExample(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {

	var err error

	// something happened
	err = fmt.Errorf("Something bad happened")

	if err != nil {
		// something didn't work out..
		errCtx := gears.NewError(c, http.StatusInternalServerError, "Sorry, just can't...")
		return errCtx
	}

	// otherwise do your thing
	w.Header().Set("Content-Type", "application/json")
	return c
}

func mainHandler(c context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("World"))
}

func main() {

	// single middleware
	http.Handle("/main", gears.NewHandler(mainHandler, middlewareExample))

	// chain middleware in the constructor
	http.Handle("/error", gears.NewHandler(mainHandler, middlewareExample, errorExample))

	// chain middleware prior using them
	withError := gears.Chain(middlewareExample, errorExample)
	http.Handle("/other_error", gears.NewHandler(mainHandler, withError))

	// ... chained middleware can be further chained.

	log.Fatal(http.ListenAndServe(":8080", nil))
}
