package controller

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestGetUsername(t *testing.T) {
	t.Run("returns empty string when no username", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest("GET", "/", nil)
		username := getUsername(r)

		assert.Empty(t, username)
	})

	t.Run("returns correct username when one given", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest("GET", "/", nil)
		r = r.WithContext(context.WithValue(r.Context(), usernameContextKey{}, "test"))
		username := getUsername(r)

		assert.Equal(t, "test", username)
	})
}

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	writeJSON(rr, map[string]interface{}{"test": true}, 503)

	assert.JSONEq(t, `{"test": true}`, rr.Body.String())
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, 503, rr.Code)
}