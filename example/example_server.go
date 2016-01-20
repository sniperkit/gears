package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zgiber/gears"
	"golang.org/x/net/context"
)

func gearHeaderExample(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {

	var err error

	// something happened
	err = fmt.Errorf("Something bad happened") // uncomment to see how errors are returned

	if err != nil {
		// something didn't work out..
		errCtx := gears.NewErrorContext(c, gears.NewStatusError(http.StatusInternalServerError, "Sorry, just can't..."))
		return errCtx
	}

	// otherwise do your thing
	w.Header().Set("Content-Type", "application/json")
	return c
}

func gearTokenExample(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {

	// otherwise do your thing
	token := r.Header.Get("Authorization")
	return context.WithValue(c, "token", token)
}

func mainHandler(c context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"message\":\"Hello World\"}"))
}

func main() {

	// single middleware
	http.Handle("/main", gears.NewHandler(mainHandler, gearHeaderExample))

	// chain middleware in the constructor
	http.Handle("/error", gears.NewHandler(mainHandler, gearTokenExample, gearHeaderExample))

	// chain middleware prior using them
	gearBox := gears.Chain(gearTokenExample, gearHeaderExample)
	http.Handle("/other_error", gears.NewHandler(mainHandler, gearBox))

	// ... chained middleware can be further chained.

	log.Fatal(http.ListenAndServe(":8080", nil))
}
