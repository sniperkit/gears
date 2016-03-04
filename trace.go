package gears

import (
	"net/http"

	"github.com/Typeform/fennel"
	"github.com/Typeform/fennel/inputs/httptrace"
)

// WithTrace returns a handler which records http request metrics (reponse time, status code, path)
// to a fennel.SimpleGatherer (metrics collection backend).
func (gear Gear) WithTrace(gatherer fennel.Gatherer) http.Handler {
	tracer := httptrace.New(gatherer)
	return tracer(gear)
}
