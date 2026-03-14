package pokemonteam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPokemonTeamSummary(t *testing.T) {
	bulbasaur := Member{
		Name:   "bulbasaur",
		Height: 7,
		Weight: 69,
		Types:  []string{"grass", "poison"},
		Stats: MemberStats{
			HP:      45,
			Attack:  49,
			Defense: 49,
			Speed:   45,
		},
	}
	pikachu := Member{
		Name:   "pikachu",
		Height: 4,
		Weight: 60,
		Types:  []string{"electric"},
		Stats: MemberStats{
			HP:      35,
			Attack:  55,
			Defense: 40,
			Speed:   90,
		},
	}
	oddish := Member{
		Name:   "oddish",
		Height: 5,
		Weight: 54,
		Types:  []string{"grass", "poison"},
		Stats: MemberStats{
			HP:      45,
			Attack:  50,
			Defense: 55,
			Speed:   30,
		},
	}

	testcases := map[string]struct {
		members  map[*Member]int
		expected Summary
	}{
		"empty team": {
			members:  map[*Member]int{},
			expected: Summary{TotalWeight: 0, AverageHeight: 0, TotalHP: 0, TypeCounts: map[string]int{}},
		},
		"team with 1 pokemon": {
			members: map[*Member]int{
				&bulbasaur: 1,
			},
			expected: Summary{TotalWeight: 69, AverageHeight: 7, TotalHP: 45, TypeCounts: map[string]int{"grass": 1, "poison": 1}},
		},
		"team with 2 pokemon": {
			members: map[*Member]int{
				&bulbasaur: 1,
				&pikachu:   1,
			},
			expected: Summary{TotalWeight: 129, AverageHeight: 5.5, TotalHP: 80, TypeCounts: map[string]int{"grass": 1, "poison": 1, "electric": 1}},
		},
		"team with duplicate types": {
			members: map[*Member]int{
				&bulbasaur: 1,
				&oddish:    1,
			},
			expected: Summary{TotalWeight: 123, AverageHeight: 6, TotalHP: 90, TypeCounts: map[string]int{"grass": 2, "poison": 2}},
		},
		"team with duplicate pokemon": {
			members: map[*Member]int{
				&bulbasaur: 2,
			},
			expected: Summary{TotalWeight: 138, AverageHeight: 7, TotalHP: 90, TypeCounts: map[string]int{"grass": 2, "poison": 2}},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			summary := GetPokemonTeamSummary(tc.members)
			assert.Equal(t, tc.expected, summary)
		})
	}
}
