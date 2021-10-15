package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./frontend", true)))
	rootGroup := router.Group("/")

	apiSvc := NewApi()
	apiSvc.Router(rootGroup)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
