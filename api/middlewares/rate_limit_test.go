package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateLimit(t *testing.T) {
	options := RateLimitOpts{Period: time.Hour, Limit: 100}

	mw := NewRateLimit(options)
	require.Equal(t, options.Period, mw.Limiter.Rate.Period)
	require.Equal(t, options.Limit, mw.Limiter.Rate.Limit)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, req)

	res := w.Result()

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "100", res.Header.Get("X-RateLimit-Limit"))
	require.Equal(t, "99", res.Header.Get("X-RateLimit-Remaining"))
	require.Equal(t, strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10), res.Header.Get("X-RateLimit-Reset"))
}
