package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/MrWestbury/terrakube-moduleregistry/services"
	"github.com/gin-gonic/gin"
)

type ModuleApi struct {
	nsSvc  services.NamespaceService
	modSvc services.ModuleService
}

func NewModuleApi(namespaceSvc services.NamespaceService, moduleSvc services.ModuleService) ModuleApi {
	newApi := &ModuleApi{
		nsSvc:  namespaceSvc,
		modSvc: moduleSvc,
	}
	return *newApi
}

func (modApi ModuleApi) Router(g *gin.RouterGroup) *gin.RouterGroup {
	apiRouter := g.Group("/:namespace/module")

	apiRouter.GET("/", modApi.ListModules)            // List namespaces
	apiRouter.POST("/", modApi.Create)                // New namespace
	apiRouter.DELETE("/:module", modApi.DeleteModule) // Delete namespace
	apiRouter.GET("/:module", modApi.GetModuleByName) // Get specific namespace

	return apiRouter
}

func (modApi ModuleApi) Create(c *gin.Context) {
	var newMod services.TerraformModule
	err := c.Bind(&newMod)
	if err != nil {
		body := ErrorResponse{
			Message: "Request body unprocessable",
		}
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, body)
		return
	}
	nsExists := modApi.nsSvc.Exists(newMod.Namespace)
	if !nsExists {
		body := ErrorResponse{
			Message: "Namespace doesn't exist",
		}
		c.AbortWithStatusJSON(http.StatusNotAcceptable, body)
		return
	}
	ns, err := modApi.nsSvc.GetNamespaceByName(newMod.Namespace)

	modExists := modApi.modSvc.Exists(*ns, newMod.Name)
	if modExists {
		body := ErrorResponse{
			Message: "Module already exists",
		}
		c.AbortWithStatusJSON(http.StatusConflict, body)
		return
	}

	createdMod, err := modApi.modSvc.CreateModule(*ns, newMod)
	if err != nil {
		body := ErrorResponse{
			Message: "Failed to create module entry",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, body)
		return
	}

	c.IndentedJSON(http.StatusCreated, createdMod)
}

func (modApi ModuleApi) ListModules(c *gin.Context) {
	ns := modApi.getNamespaceFromRequest(c)

	moduleList := modApi.modSvc.ListModules(*ns)
	response := ResponseListModules{
		Modules: *moduleList,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (modApi ModuleApi) GetModuleByName(c *gin.Context) {
	ns := modApi.getNamespaceFromRequest(c)

	moduleName := c.Param("module")
	module, err := modApi.modSvc.GetModuleByName(*ns, moduleName)
	if err != nil {
		if err == services.ErrModuleNotFound {
			response := ErrorResponse{
				Message: "Module not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return
		}
		log.Error("error getting module", err)
		response := ErrorResponse{
			Message: "Unknown error getting module",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	c.IndentedJSON(http.StatusOK, module)
}

func (modApi ModuleApi) DeleteModule(c *gin.Context) {
	ns := modApi.getNamespaceFromRequest(c)

	moduleName := c.Param("module")
	modApi.modSvc.DeleteModule(*ns, moduleName)
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (modApi ModuleApi) getNamespaceFromRequest(c *gin.Context) *services.Namespace {
	namespaceName := c.Param("namespace")
	ns, err := modApi.nsSvc.GetNamespaceByName(namespaceName)
	if err != nil {
		if err == services.ErrNamespaceNotFound {
			response := ErrorResponse{
				Message: "Namespace not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return nil
		}

		response := ErrorResponse{
			Message: "Unknown error listing modules",
		}

		log.Error("error getting namespace", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return nil
	}

	return ns
}
