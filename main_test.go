package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPokemonTeam(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	testcases := map[string]struct {
		query  string
		status int
		body   string
	}{
		"no query": {
			query:  "",
			status: http.StatusBadRequest,
			body:   "{\"message\":\"At least 1 pokemon name is required\"}",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/pokemon/team?names="+tc.query, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.status, w.Code)
			assert.Equal(t, tc.body, w.Body.String())
		})
	}
}
