package pokemonteam

import (
	"fmt"
	"strings"

	pokego "github.com/JoshGuarino/PokeGo/pkg"
)

type Member struct {
	Name   string      `json:"name"`
	Height int         `json:"height"`
	Weight int         `json:"weight"`
	Types  []string    `json:"types"`
	Stats  MemberStats `json:"stats"`
	Image  string      `json:"image"`
}

type MemberStats struct {
	HP      int `json:"hp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
}

type Summary struct {
	TotalWeight   int            `json:"total_weight"`
	AverageHeight float64        `json:"average_height"`
	TotalHP       int            `json:"total_hp"`
	TypeCounts    map[string]int `json:"type_counts"`
}

const MaxPokemon = 6

type PokemonAPIService struct {
	pokemonapi pokego.PokeGo
}

func NewPokemonAPIService(pokemonapi pokego.PokeGo) Service {
	return &PokemonAPIService{pokemonapi: pokemonapi}
}

type Service interface {
	GetMembers(names []string) ([]Member, error)
	GetSummary(team map[*Member]int) Summary
}

// Assumes no duplicates
func (s *PokemonAPIService) GetMembers(names []string) ([]Member, error) {
	result := make([]Member, 0, len(names))

	for _, name := range names {
		pokemon, err := s.pokemonapi.Pokemon.GetPokemon(name)
		if err != nil {
			if strings.Contains(err.Error(), "invalid") {
				return nil, fmt.Errorf("%s is not a valid pokemon name", name)
			}

			return nil, fmt.Errorf("error fetching data for pokemon %s: %w", name, err)
		}
		member := Member{
			Name:   pokemon.Name,
			Height: pokemon.Height,
			Weight: pokemon.Weight,
			Types:  make([]string, 0, len(pokemon.Types)),
			Stats:  MemberStats{},
			Image:  pokemon.Sprites.FrontDefault,
		}
		for _, stat := range pokemon.Stats {
			statName := stat.Stat.Name
			switch statName {
			case "hp":
				member.Stats.HP = stat.BaseStat
			case "attack":
				member.Stats.Attack = stat.BaseStat
			case "defense":
				member.Stats.Defense = stat.BaseStat
			case "speed":
				member.Stats.Speed = stat.BaseStat
			}
		}
		for _, t := range pokemon.Types {
			member.Types = append(member.Types, t.Type.Name)
		}
		result = append(result, member)
	}

	return result, nil
}

func (s *PokemonAPIService) GetSummary(team map[*Member]int) Summary {
	return GetPokemonTeamSummary(team)
}

func GetPokemonTeamSummary(team map[*Member]int) Summary {
	summary := Summary{
		TypeCounts: make(map[string]int),
	}
	if len(team) == 0 {
		return summary
	}
	totalHeight := 0
	for member, count := range team {
		summary.TotalWeight += member.Weight * count
		totalHeight += member.Height * count
		summary.TotalHP += member.Stats.HP * count
		for _, t := range member.Types {
			summary.TypeCounts[t] += count
		}
	}
	totalTeamCount := 0
	for _, count := range team {
		totalTeamCount += count
	}
	summary.AverageHeight = float64(totalHeight) / float64(totalTeamCount)
	return summary
}
