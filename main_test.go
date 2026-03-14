package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPokemonTeam(t *testing.T) {
	router := setupRouter()

	testcases := map[string]struct {
		query  string
		status int
		body   string
	}{
		"empty query": {
			query:  "",
			status: http.StatusBadRequest,
			body:   "{\"message\":\"At least 1 pokemon name is required\"}",
		},
		"empty query with commas": {
			query:  ",,,,",
			status: http.StatusBadRequest,
			body:   "{\"message\":\"At least 1 pokemon name is required\"}",
		},
		"empty query with whitespace": {
			query:  "  ,  ,  ",
			status: http.StatusBadRequest,
			body:   "{\"message\":\"At least 1 pokemon name is required\"}",
		},
		"too many pokemon": {
			query:  "pikachu,bulbasaur,charmander,squirtle,eevee,meowth,psyduck",
			status: http.StatusBadRequest,
			body:   "{\"message\":\"No more than 6 pokemon names are allowed\"}",
		},
	}

	for name, tc := range testcases {
		w := httptest.NewRecorder()
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/pokemon/team?names="+tc.query, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.status, w.Code)
			assert.Equal(t, tc.body, w.Body.String())
		})
	}
}
