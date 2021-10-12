package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/MrWestbury/terrakube-moduleregistry/services"
	"github.com/gin-gonic/gin"
)

type SystemApi struct {
	nsSvc  services.NamespaceService
	modSvc services.ModuleService
	sysSvc services.SystemService
}

func NewSystemApi(namespaceSvc services.NamespaceService, moduleSvc services.ModuleService, systemSvc services.SystemService) *SystemApi {
	api := &SystemApi{
		nsSvc:  namespaceSvc,
		modSvc: moduleSvc,
		sysSvc: systemSvc,
	}

	return api
}

func (sysApi SystemApi) Router(g *gin.RouterGroup) *gin.RouterGroup {
	apiRouter := g.Group("/:module/system")

	apiRouter.GET("/", sysApi.ListSystems)      // List namespaces
	apiRouter.POST("/", sysApi.CreateSystem)    // New namespace
	apiRouter.DELETE("/:system", sysApi.Dummy)  // Delete namespace
	apiRouter.GET("/:system", sysApi.GetSystem) // Get specific namespace

	return apiRouter
}

func (sysApi SystemApi) Dummy(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (sysApi SystemApi) CreateSystem(c *gin.Context) {
	var newSys services.TerraformSystem
	err := c.Bind(&newSys)
	if err != nil {
		body := ErrorResponse{
			Message: "Request body is unprocessable",
		}
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, body)
	}

	namespaceName := c.Param("namespace")
	ns, err := sysApi.nsSvc.GetNamespaceByName(namespaceName)
	if err != nil {
		if err == services.ErrNamespaceNotFound {
			body := ErrorResponse{
				Message: "Namespace not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	newSys.Namespace = ns.Name

	moduleName := c.Param("module")
	mod, err := sysApi.modSvc.GetModuleByName(*ns, moduleName)
	if err != nil {
		if err == services.ErrModuleNotFound {
			body := ErrorResponse{
				Message: "Module not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	newSys.Module = mod.Name

	system, err := sysApi.sysSvc.CreateSystem(newSys)
	if err != nil {
		if err == services.ErrSystemAlreadyExists {
			response := ErrorResponse{
				Message: "System already exists",
			}
			c.AbortWithStatusJSON(http.StatusConflict, response)
		}
		response := ErrorResponse{
			Message: "Error creating system",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
	}

	c.IndentedJSON(http.StatusCreated, system)
}

func (sysApi SystemApi) ListSystems(c *gin.Context) {
	namespaceName := c.Param("namespace")
	ns, err := sysApi.nsSvc.GetNamespaceByName(namespaceName)
	if err != nil {
		if err == services.ErrNamespaceNotFound {
			body := ErrorResponse{
				Message: "Namespace not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	moduleName := c.Param("module")
	mod, err := sysApi.modSvc.GetModuleByName(*ns, moduleName)
	if err != nil {
		if err == services.ErrModuleNotFound {
			body := ErrorResponse{
				Message: "Module not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	systemList := sysApi.sysSvc.ListSystemsByModule(*mod)

	response := ResponseListSystems{
		Namespace: mod.Namespace,
		Module:    mod.Name,
		Systems:   *systemList,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (sysApi SystemApi) GetSystem(c *gin.Context) {
	namespaceName := c.Param("namespace")
	ns, err := sysApi.nsSvc.GetNamespaceByName(namespaceName)
	if err != nil {
		if err == services.ErrNamespaceNotFound {
			body := ErrorResponse{
				Message: "Namespace not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	moduleName := c.Param("module")
	mod, err := sysApi.modSvc.GetModuleByName(*ns, moduleName)
	if err != nil {
		if err == services.ErrModuleNotFound {
			body := ErrorResponse{
				Message: "Module not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	systemName := c.Param("system")
	system, err := sysApi.sysSvc.GetSystemByName(*mod, systemName)
	if err != nil {
		if err == services.ErrSystemNotFound {
			body := ErrorResponse{
				Message: "System not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
		log.Errorf("error getting system: %v", err)
		body := ErrorResponse{
			Message: "Unknown error getting system",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, body)
	}

	c.IndentedJSON(http.StatusOK, system)
}
