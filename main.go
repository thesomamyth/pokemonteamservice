package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/pokemon/team", getPokemonTeamHandler)
	return router
}

const maxPokemon = 6

func getPokemonTeamHandler(c *gin.Context) {
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

	teamData, err := getPokemonTeamData(teamNames.UniqueNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching pokemon data"})
		return
	}

	c.JSON(http.StatusNotImplemented, gin.H{"team": teamData})
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

	for i := range splitNames {
		trimmed := strings.TrimSpace(splitNames[i])
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

// Assumes no duplicates
func getPokemonTeamData(names []string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}
