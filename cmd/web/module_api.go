package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/MrWestbury/terraxen/services"
	"github.com/gin-gonic/gin"
)

type ModuleApi struct {
	nsSvc  services.NamespaceService
	modSvc services.ModuleService
	sysSvc services.SystemService
}

func NewModuleApi(namespaceSvc services.NamespaceService, moduleSvc services.ModuleService, systemSvc services.SystemService) ModuleApi {
	newApi := &ModuleApi{
		nsSvc:  namespaceSvc,
		modSvc: moduleSvc,
		sysSvc: systemSvc,
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
	var moduleRequest RequestNewModule
	err := c.Bind(&moduleRequest)
	if err != nil {
		body := ErrorResponse{
			Message: "Request body unprocessable",
		}
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, body)
		return
	}
	ns := modApi.getNamespaceFromRequest(c)
	if ns == nil {
		return
	}

	newModule := services.NewTerraformModule{
		Name:      moduleRequest.Name,
		Namespace: ns.Name,
	}

	createdMod, err := modApi.modSvc.CreateModule(newModule)
	if err != nil {
		if err == services.ErrModuleAlreadyExists {
			response := ErrorResponse{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}
		response := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create module entry",
		}
		c.AbortWithStatusJSON(response.Code, response)
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
	module := modApi.getModuleFromRequest(c)

	systemList := modApi.sysSvc.ListSystemsByModule(*module)
	if len(*systemList) > 0 {
		response := ErrorResponse{
			Message: "Module has dependant systems",
		}
		c.AbortWithStatusJSON(http.StatusConflict, response)
		return
	}

	modApi.modSvc.DeleteModule(module)
	c.AbortWithStatus(http.StatusOK)
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
			Message: "Unknown error getting namespace",
		}

		log.Error("error getting namespace", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return nil
	}

	return ns
}

func (modApi ModuleApi) getModuleFromRequest(c *gin.Context) *services.TerraformModule {
	ns := modApi.getNamespaceFromRequest(c)
	if ns == nil {
		return nil
	}

	moduleName := c.Param("module")
	module, err := modApi.modSvc.GetModuleByName(*ns, moduleName)
	if err != nil {
		if err == services.ErrModuleNotFound {
			response := ErrorResponse{
				Message: "Module not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return nil
		}
		response := ErrorResponse{
			Message: "Unknown error getting module",
		}

		log.Error("error getting module", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return nil
	}

	return module
}
