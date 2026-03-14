package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/pokemon/team", getPokemonTeam)
	return router
}

func getPokemonTeam(c *gin.Context) {
	names, exists := c.GetQuery("names")
	if !exists || names == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "At least 1 pokemon name is required"})
		return
	}

	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}
