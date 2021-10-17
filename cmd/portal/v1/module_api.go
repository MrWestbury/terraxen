package v1

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/MrWestbury/terraxen/services"
	"github.com/gin-gonic/gin"
)

type ModuleApi struct {
	helper ApiHelper
}

func NewModuleApi(helper ApiHelper) ModuleApi {
	newApi := &ModuleApi{
		helper: helper,
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
			Code:    http.StatusUnprocessableEntity,
			Message: "Request body unprocessable",
		}
		c.AbortWithStatusJSON(body.Code, body)
		return
	}
	ns := modApi.helper.GetNamespaceFromRequest(c)
	if ns == nil {
		return
	}

	newModule := services.NewTerraformModule{
		Name:      moduleRequest.Name,
		Namespace: ns.Name,
	}

	createdMod, err := modApi.helper.ModuleSvc.CreateModule(newModule)
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
	ns := modApi.helper.GetNamespaceFromRequest(c)

	moduleList := modApi.helper.ModuleSvc.ListModules(*ns)
	response := ResponseListModules{
		Modules: *moduleList,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (modApi ModuleApi) GetModuleByName(c *gin.Context) {
	ns := modApi.helper.GetNamespaceFromRequest(c)

	moduleName := c.Param("module")
	module, err := modApi.helper.ModuleSvc.GetModuleByName(*ns, moduleName)
	if err != nil {
		if err == services.ErrModuleNotFound {
			response := ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Module not found",
			}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}
		log.Error("error getting module", err)
		response := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unknown error getting module",
		}
		c.AbortWithStatusJSON(response.Code, response)
		return
	}

	c.IndentedJSON(http.StatusOK, module)
}

func (modApi ModuleApi) DeleteModule(c *gin.Context) {
	module := modApi.helper.GetModuleFromRequest(c)
	if module == nil {
		return
	}

	systemList := modApi.helper.SystemSvc.ListSystemsByModule(*module)
	if len(*systemList) > 0 {
		response := ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Module has dependant systems",
		}
		c.AbortWithStatusJSON(response.Code, response)
		return
	}

	modApi.helper.ModuleSvc.DeleteModule(*module)
	c.Status(http.StatusOK)
}
