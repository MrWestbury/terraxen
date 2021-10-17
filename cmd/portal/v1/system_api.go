package v1

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/MrWestbury/terraxen/services"
	"github.com/gin-gonic/gin"
)

type SystemApi struct {
	helper ApiHelper
}

func NewSystemApi(helper ApiHelper) *SystemApi {
	api := &SystemApi{
		helper: helper,
	}

	return api
}

func (sysApi SystemApi) Router(g *gin.RouterGroup) *gin.RouterGroup {
	apiRouter := g.Group("/:module/system")

	apiRouter.GET("/", sysApi.ListSystems)            // List namespaces
	apiRouter.POST("/", sysApi.CreateSystem)          // New namespace
	apiRouter.DELETE("/:system", sysApi.DeleteSystem) // Delete namespace
	apiRouter.GET("/:system", sysApi.GetSystem)       // Get specific namespace

	return apiRouter
}

func (sysApi SystemApi) CreateSystem(c *gin.Context) {
	var newSys RequestNewSystem
	err := c.Bind(&newSys)
	if err != nil {
		body := ErrorResponse{
			Message: "Request body is unprocessable",
		}
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, body)
		return
	}

	module := sysApi.helper.GetModuleFromRequest(c)
	if module == nil {
		return
	}

	exists := sysApi.helper.SystemSvc.ExistsByName(*module, newSys.Name)
	if exists {
		response := ErrorResponse{
			Code:    http.StatusConflict,
			Message: "System already exists",
		}
		c.AbortWithStatusJSON(response.Code, response)
		return
	}

	newSystem := services.NewTerraformSystem{
		Name:   newSys.Name,
		Module: *module,
	}

	system, err := sysApi.helper.SystemSvc.CreateSystem(newSystem)
	if err != nil {
		response := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "System creation failed",
		}
		c.AbortWithStatusJSON(response.Code, response)
		return
	}

	c.IndentedJSON(http.StatusCreated, system)
}

func (sysApi SystemApi) ListSystems(c *gin.Context) {
	mod := sysApi.helper.GetModuleFromRequest(c)
	if mod == nil {
		return
	}

	systemList := sysApi.helper.SystemSvc.ListSystemsByModule(*mod)

	response := ResponseListSystems{
		Namespace: mod.Namespace,
		Module:    mod.Name,
		Systems:   *systemList,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (sysApi SystemApi) GetSystem(c *gin.Context) {
	mod := sysApi.helper.GetModuleFromRequest(c)
	if mod == nil {
		return
	}

	systemName := c.Param("system")
	system, err := sysApi.helper.SystemSvc.GetSystemByName(*mod, systemName)
	if err != nil {
		if err == services.ErrSystemNotFound {
			response := ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "System not found",
			}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}
		log.Errorf("error getting system: %v", err)
		response := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unknown error getting system",
		}
		c.AbortWithStatusJSON(response.Code, response)
		return
	}

	c.IndentedJSON(http.StatusOK, system)
}

func (sysApi SystemApi) DeleteSystem(c *gin.Context) {
	sys := sysApi.helper.GetSystemFromRequest(c)
	if sys == nil {
		return
	}

	hasChildren := sysApi.helper.VersionSvc.HasChildren(*sys)
	if hasChildren {
		response := ErrorResponse{
			Code:    http.StatusConflict,
			Message: "System has child versions. Delete the versions before proceeding",
		}
		c.AbortWithStatusJSON(response.Code, response)
		return
	}

	err := sysApi.helper.SystemSvc.Delete(*sys)
	if err != nil {
		log.Printf("Error: %v", err)
		response := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete system",
		}
		c.AbortWithStatusJSON(response.Code, response)
	}

	c.Status(http.StatusOK)
}
