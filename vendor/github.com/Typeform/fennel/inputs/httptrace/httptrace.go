package httptrace

import (
	"log"
	"net/http"
	"time"

	"github.com/Typeform/fennel"
)

var (
	httpMetric = fennel.NewMetric("http_trace", "ms")
)

type httpWriter struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
}

func (w *httpWriter) WriteHeader(status int) {
	if !w.wroteHeader {
		w.status = status
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w httpWriter) Status() int {
	return w.status
}

// New takes a gatherer and returns a function which has the common
// middleware signature of func(h http.Handler) http.Handler
// this middleware wraps a http.Handler and records the following
// stats to the Gatherer backend:
// - request status
// - request path
// - request response time
func New(g fennel.Gatherer) func(h http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			hw := &httpWriter{ResponseWriter: w}

			t1 := time.Now()
			h.ServeHTTP(hw, r)

			status := hw.Status()
			if status == 0 {
				// certain responses might write the header without calling the
				// custom httpWriter (?), resulting in missing response code from
				// the hw.Status() while the response has a status code already..
				status = 200
			}

			go func(status int, path string) {
				// create new datapoint
				tags := map[string]interface{}{
					"status": status,
					"method": r.Method,
				}

				datapoint := fennel.NewDatapoint(httpMetric, time.Since(t1), time.Now().UTC(), tags)
				err := g.Gather(datapoint)
				if err != nil {
					log.Println(err)
				}
			}(status, r.URL.Path)
		}

		return http.HandlerFunc(fn)
	}
}
