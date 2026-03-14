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
		"valid query with 1 pokemon": {
			query:  "bulbasaur",
			status: http.StatusOK,
			body:   "{\"team\":{\"members\":[{\"name\":\"bulbasaur\",\"height\":7,\"weight\":69,\"types\":[\"grass\",\"poison\"],\"stats\":{\"hp\":45,\"attack\":49,\"defense\":49,\"speed\":45},\"image\":\"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png\"}],\"summary\":{\"total_weight\":69,\"average_height\":7,\"total_hp\":45,\"type_counts\":{\"grass\":1,\"poison\":1}}}}",
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

func TestGetPokemonTeamSummary(t *testing.T) {
	testcases := map[string]struct {
		members  []PokemonTeamMember
		expected PokemonTeamSummary
	}{
		"empty team": {
			members:  []PokemonTeamMember{},
			expected: PokemonTeamSummary{TotalWeight: 0, AverageHeight: 0, TotalHP: 0, TypeCounts: map[string]int{}},
		},
		"team with 1 pokemon": {
			members: []PokemonTeamMember{
				{Name: "bulbasaur", Height: 7, Weight: 69, Types: []string{"grass", "poison"}, Stats: PokemonTeamMemberStats{HP: 45}},
			},
			expected: PokemonTeamSummary{TotalWeight: 69, AverageHeight: 7, TotalHP: 45, TypeCounts: map[string]int{"grass": 1, "poison": 1}},
		},
		"team with 2 pokemon": {
			members: []PokemonTeamMember{
				{Name: "bulbasaur", Height: 7, Weight: 69, Types: []string{"grass", "poison"}, Stats: PokemonTeamMemberStats{HP: 45}},
				{Name: "pikachu", Height: 4, Weight: 60, Types: []string{"electric"}, Stats: PokemonTeamMemberStats{HP: 35}},
			},
			expected: PokemonTeamSummary{TotalWeight: 129, AverageHeight: 5.5, TotalHP: 80, TypeCounts: map[string]int{"grass": 1, "poison": 1, "electric": 1}},
		},
		"team with duplicate types": {
			members: []PokemonTeamMember{
				{Name: "bulbasaur", Height: 7, Weight: 69, Types: []string{"grass", "poison"}, Stats: PokemonTeamMemberStats{HP: 45}},
				{Name: "oddish", Height: 5, Weight: 54, Types: []string{"grass", "poison"}, Stats: PokemonTeamMemberStats{HP: 45}},
			},
			expected: PokemonTeamSummary{TotalWeight: 123, AverageHeight: 6, TotalHP: 90, TypeCounts: map[string]int{"grass": 2, "poison": 2}},
		},
		"team with duplicate pokemon": {
			members: []PokemonTeamMember{
				{Name: "bulbasaur", Height: 7, Weight: 69, Types: []string{"grass", "poison"}, Stats: PokemonTeamMemberStats{HP: 45}},
				{Name: "bulbasaur", Height: 7, Weight: 69, Types: []string{"grass", "poison"}, Stats: PokemonTeamMemberStats{HP: 45}},
			},
			expected: PokemonTeamSummary{TotalWeight: 138, AverageHeight: 7, TotalHP: 90, TypeCounts: map[string]int{"grass": 2, "poison": 2}},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			summary := getPokemonTeamSummary(tc.members)
			assert.Equal(t, tc.expected, summary)
		})
	}
}
