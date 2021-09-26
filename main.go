package main

import (
	"github.com/MrWestbury/terrakube-moduleregistry/services"
	"net/http"
	"os"

	"github.com/MrWestbury/terrakube-moduleregistry/api"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var modulesRootPath string = "/modules/v1"

func main() {
	log.Info("Module registry starting")

	svcOpts := services.Options{
		Hostname: os.Getenv("MONGO_HOSTNAME"),
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Database: os.Getenv("MONGO_DB"),
		ExtraOpts: map[string]string{
			"retryWrites": "true",
			"w":           "majority",
		},
	}

	nsSvc := services.NewNamespaceService(svcOpts)
	modSvc := services.NewModuleService(svcOpts)
	sysSvc := services.NewSystemService(svcOpts)
	verSvc := services.NewVersionService(svcOpts)

	router := gin.Default()
	rootGroup := router.Group("/")

	router.GET("/.well-known/terraform.json", getManifest)

	moduleRegistry := api.NewModuleRegistry(*nsSvc, *modSvc, *sysSvc, *verSvc)
	moduleRegistry.Router(rootGroup)

	apiSvc := api.Api{}
	apiSvc.Router(rootGroup)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

func getManifest(c *gin.Context) {
	data := map[string]string{
		"modules.v1": modulesRootPath,
	}
	c.IndentedJSON(http.StatusOK, data)
}
