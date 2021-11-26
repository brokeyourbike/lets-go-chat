package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if possible.
func Recoverer(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
					logger.WithFields(logrus.Fields{
						"stacktrace": string(debug.Stack()),
					}).Warnf("%+v", rvr)

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
