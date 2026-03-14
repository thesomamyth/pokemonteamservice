package main

import (
	"fmt"
	"net/http"
	"strings"

	pokego "github.com/JoshGuarino/PokeGo/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

type service struct {
	r          *gin.Engine
	pokemonapi pokego.PokeGo // TODO: extract to interface
}

func setupRouter() *gin.Engine {
	pokemonClient := pokego.NewClient()
	s := service{
		pokemonapi: pokemonClient,
	}
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/pokemon/team", s.getPokemonTeamHandler)
	return router
}

const maxPokemon = 6

func (s service) getPokemonTeamHandler(c *gin.Context) {
	query, exists := c.GetQuery("names")
	if !exists || query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "At least 1 pokemon name is required"})
		return
	}

	teamNames := filterNames(query)

	if len(teamNames.Names) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "At least 1 pokemon name is required"})
		return
	}

	if len(teamNames.Names) > maxPokemon {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("No more than %d pokemon names are allowed", maxPokemon)})
		return
	}

	members, err := s.getPokemonTeamMembers(teamNames.UniqueNames)
	if err != nil {
		if strings.Contains(err.Error(), "is not a valid pokemon name") {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		// Since this is not sensitive data, we bubble up the error message
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	fullTeam := make(map[*PokemonTeamMember]int)
	for i := range members {
		fullTeam[&members[i]] = teamNames.NameCounts[members[i].Name]
	}

	summary := getPokemonTeamSummary(fullTeam)

	response := PokemonTeamResponse{
		Members: members,
		Summary: summary,
	}

	c.JSON(http.StatusOK, gin.H{"team": response})
}

type teamNameInfo struct {
	Names       []string
	UniqueNames []string
	NameCounts  map[string]int
}

func filterNames(names string) teamNameInfo {
	splitNames := strings.Split(names, ",")
	result := teamNameInfo{
		Names:       make([]string, 0, len(splitNames)),
		UniqueNames: make([]string, 0, len(splitNames)),
		NameCounts:  make(map[string]int, len(splitNames)),
	}

	for _, splitName := range splitNames {
		trimmed := strings.TrimSpace(splitName)
		if trimmed == "" {
			continue
		}

		lowercaseName := strings.ToLower(trimmed)
		result.Names = append(result.Names, lowercaseName)
		result.NameCounts[lowercaseName]++
		if result.NameCounts[lowercaseName] == 1 {
			result.UniqueNames = append(result.UniqueNames, lowercaseName)
		}
	}
	return result
}

type PokemonTeamResponse struct {
	Members []PokemonTeamMember `json:"members"`
	Summary PokemonTeamSummary  `json:"summary"`
}

type PokemonTeamMember struct {
	Name   string                 `json:"name"`
	Height int                    `json:"height"`
	Weight int                    `json:"weight"`
	Types  []string               `json:"types"`
	Stats  PokemonTeamMemberStats `json:"stats"`
	Image  string                 `json:"image"`
}

type PokemonTeamMemberStats struct {
	HP      int `json:"hp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
}

type PokemonTeamSummary struct {
	TotalWeight   int            `json:"total_weight"`
	AverageHeight float64        `json:"average_height"`
	TotalHP       int            `json:"total_hp"`
	TypeCounts    map[string]int `json:"type_counts"`
}

// Assumes no duplicates
func (s service) getPokemonTeamMembers(names []string) ([]PokemonTeamMember, error) {
	result := make([]PokemonTeamMember, 0, len(names))

	for _, name := range names {
		pokemon, err := s.pokemonapi.Pokemon.GetPokemon(name)
		if err != nil {
			if strings.Contains(err.Error(), "invalid") {
				return nil, fmt.Errorf("%s is not a valid pokemon name", name)
			}

			return nil, fmt.Errorf("error fetching data for pokemon %s: %w", name, err)
		}
		member := PokemonTeamMember{
			Name:   pokemon.Name,
			Height: pokemon.Height,
			Weight: pokemon.Weight,
			Types:  make([]string, 0, len(pokemon.Types)),
			Stats:  PokemonTeamMemberStats{},
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

func getPokemonTeamSummary(team map[*PokemonTeamMember]int) PokemonTeamSummary {
	summary := PokemonTeamSummary{
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
