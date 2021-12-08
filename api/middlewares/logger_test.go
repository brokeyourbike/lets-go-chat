package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	hook := test.NewGlobal()
	r := chi.NewRouter()

	r.Use(Logger)
	r.Get("/", testHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	assert.Equal(t, 2, len(hook.Entries))
	assert.Contains(t, hook.Entries[0].Message, "Request:")
	assert.Contains(t, hook.Entries[1].Message, "Response:")
}
