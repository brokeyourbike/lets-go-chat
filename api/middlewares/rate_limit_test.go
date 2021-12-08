package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRateLimit(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	rl := NewRateLimit(RateLimitOpts{Period: time.Hour, Limit: 100})
	require.Equal(t, time.Hour, rl.Limiter.Rate.Period)
	require.Equal(t, int64(100), rl.Limiter.Rate.Limit)

	r := chi.NewRouter()
	r.Get("/", rl.Handle(testHandler))

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "100", res.Header.Get("X-RateLimit-Limit"))
	require.Equal(t, "99", res.Header.Get("X-RateLimit-Remaining"))
	require.Equal(t, strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10), res.Header.Get("X-RateLimit-Reset"))
}
