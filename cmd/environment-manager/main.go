package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		resp := map[string]interface{}{
			"message": "ok",
		}
		c.IndentedJSON(http.StatusOK, resp)
	})

	err := router.Run(":8000")
	if err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
