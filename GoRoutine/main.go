package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.Run("localhost:8080")
}

type album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  string `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "GNAT", Artist: "Eminem", Price: "$9.99"},
	{ID: "2", Title: "Theme Alfred", Artist: "Eminem", Price: "9.99"},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
