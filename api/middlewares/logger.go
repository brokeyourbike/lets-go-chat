package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logRequestFields(r)).Infof("Request: %s %s", r.Method, r.URL.Path)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			log.WithFields(log.Fields{
				"status":  ww.Status(),
				"bytes":   ww.BytesWritten(),
				"elapsed": float64(time.Since(t1).Nanoseconds()) / 1000000.0, // in milliseconds
			}).Logf(statusLevel(ww.Status()), "Response: %d %s", ww.Status(), statusLabel(ww.Status()))
		}()

		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}

func logRequestFields(r *http.Request) log.Fields {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	requestURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	return log.Fields{
		"requestURL":    requestURL,
		"requestMethod": r.Method,
		"requestPath":   r.URL.Path,
		"remoteIP":      r.RemoteAddr,
		"proto":         r.Proto,
	}
}

func statusLevel(status int) log.Level {
	switch {
	case status <= 0:
		return log.WarnLevel
	case status < 400:
		return log.InfoLevel
	case status >= 400 && status < 500:
		return log.WarnLevel
	case status >= 500:
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}

func statusLabel(status int) string {
	switch {
	case status >= 100 && status < 300:
		return "OK"
	case status >= 300 && status < 400:
		return "Redirect"
	case status >= 400 && status < 500:
		return "Client Error"
	case status >= 500:
		return "Server Error"
	default:
		return "Unknown"
	}
}
