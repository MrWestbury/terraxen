package main

import (
	"net/http"

	"github.com/MrWestbury/terrakube-moduleregistry/api"
	"github.com/MrWestbury/terrakube-moduleregistry/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var modulesRootPath string = "/modules/v1"

var data DataController

func main() {
	log.Info("Module registry starting")

	router := gin.Default()
	router.GET("/.well-known/terraform.json", getManifest)

	moduleGroup := router.Group(modulesRootPath)
	moduleGroup.GET(":namespace/:name/:provider/versions", getModuleVersions)

	api := api.Api{}
	api.Router(router.Group("/"))

	data = DataController{}

	router.Run("localhost:8080")
}

func getManifest(c *gin.Context) {
	data := map[string]string{
		"modules.v1": modulesRootPath,
	}
	c.IndentedJSON(http.StatusOK, data)
}

func getModuleVersions(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	provider := c.Param("provider")

	versions := data.GetModuleVersions(namespace, name, provider)

	terraformModule := models.TerraformModule{}

	for _, version := range versions {
		newVersion := models.ModuleVersion{
			Version: version,
		}

		terraformModule.Versions = append(terraformModule.Versions, newVersion)
	}

	parent := models.ModuleVersionsParent{
		Modules: []models.TerraformModule{terraformModule},
	}

	c.IndentedJSON(200, parent)
}
