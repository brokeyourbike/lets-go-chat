package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorLogger(t *testing.T) {
	hook := test.NewGlobal()

	mw := ErrorLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Foo error", http.StatusInternalServerError)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "Foo error\n", hook.LastEntry().Message)
}
