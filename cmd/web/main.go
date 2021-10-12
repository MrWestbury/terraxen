package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()
	rootGroup := router.Group("/")

	apiSvc := NewApi()
	apiSvc.Router(rootGroup)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
