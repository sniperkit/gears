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
	// err = fmt.Errorf("Something bad happened") // uncomment to see how errors are returned

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
	w.Header().Set("Authorization", token) // just so we see something in the result..
	return context.WithValue(c, "token", token)
}

func gearErrorExample(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {

	// let's say we have an error here
	err := fmt.Errorf("Can't do that, sorry... would you like a hot beverage?")
	return gears.NewErrorContext(c, gears.NewStatusError(http.StatusInternalServerError, err.Error()))
}

func mainHandler(c context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"message\":\"Hello World\"}"))
}

func main() {

	base := gears.NewHandler(mainHandler, gearHeaderExample)

	// single middleware
	http.Handle("/main", base)

	// extend base handler with a new gear (middleware)
	http.Handle("/main/extended", gears.NewHandler(base, gearTokenExample))

	// chain middleware in the handler constructor...
	http.Handle("/error", gears.NewHandler(mainHandler, gearErrorExample, gearHeaderExample))

	//...or chain middleware before using them in the constructor
	gearBox := gears.Chain(gearTokenExample, gearHeaderExample)
	http.Handle("/gearbox", gears.NewHandler(mainHandler, gearBox))

	// tip: chained middleware can be chained further.

	log.Fatal(http.ListenAndServe(":8080", nil))
}
