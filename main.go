package main

import (
	"fmt"
	"net/http"
	"pokemonteamservice/pokemonteam"
	"strings"

	pokego "github.com/JoshGuarino/PokeGo/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

type pokemonTeamRouter struct {
	teamService pokemonteam.Service
}

func setupRouter() *gin.Engine {
	s := pokemonteam.NewPokemonAPIService(pokego.NewClient())
	r := &pokemonTeamRouter{teamService: s}

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/pokemon/team", r.getPokemonTeamHandler)
	return router
}

type PokemonTeamResponse struct {
	Members []pokemonteam.Member `json:"members"`
	Summary pokemonteam.Summary  `json:"summary"`
}

func (r *pokemonTeamRouter) getPokemonTeamHandler(c *gin.Context) {
	query, exists := c.GetQuery("names")
	if !exists || query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least 1 pokemon name is required"})
		return
	}

	teamNames := filterNames(query)
	if len(teamNames.Names) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least 1 pokemon name is required"})
		return
	}

	if len(teamNames.Names) > pokemonteam.MaxPokemon {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No more than %d pokemon names are allowed", pokemonteam.MaxPokemon)})
		return
	}

	members, err := r.teamService.GetMembers(teamNames.UniqueNames)
	if err != nil {
		if strings.Contains(err.Error(), "is not a valid pokemon name") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Since this is not sensitive data, we can bubble up the error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fullTeam := make(map[*pokemonteam.Member]int)
	for i := range members {
		fullTeam[&members[i]] = teamNames.NameCounts[members[i].Name]
	}

	summary := r.teamService.GetSummary(fullTeam)

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
