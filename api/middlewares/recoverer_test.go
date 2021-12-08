package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecoverer(t *testing.T) {
	panicingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("foo")
	})

	hook := test.NewGlobal()
	r := chi.NewRouter()

	r.Use(Recoverer)
	r.Get("/", panicingHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)

	_, exists := hook.LastEntry().Data["stacktrace"]
	assert.True(t, exists)
}
