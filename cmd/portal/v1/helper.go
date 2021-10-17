package v1

import (
	"github.com/MrWestbury/terraxen/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ApiHelper struct {
	NamespaceSvc services.NamespaceService
	ModuleSvc    services.ModuleService
	SystemSvc    services.SystemService
	VersionSvc   services.VersionService
}

func NewHelper(nsSvc services.NamespaceService, modSvc services.ModuleService, sysSvc services.SystemService, verSvc services.VersionService) *ApiHelper {
	helper := &ApiHelper{
		NamespaceSvc: nsSvc,
		ModuleSvc:    modSvc,
		SystemSvc:    sysSvc,
		VersionSvc:   verSvc,
	}

	return helper
}

func (helper ApiHelper) GetNamespaceFromRequest(c *gin.Context) *services.TerraformNamespace {
	namespaceName := c.Param("namespace")

	ns, err := helper.NamespaceSvc.GetNamespaceByName(namespaceName)
	if err != nil {
		if err == services.ErrNamespaceNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return nil
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	return ns
}

func (helper ApiHelper) GetModuleFromRequest(c *gin.Context) *services.TerraformModule {
	ns := helper.GetNamespaceFromRequest(c)
	if ns == nil {
		return nil
	}

	moduleName := c.Param("module")
	module, err := helper.ModuleSvc.GetModuleByName(*ns, moduleName)
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

func (helper ApiHelper) GetSystemFromRequest(c *gin.Context) *services.TerraformSystem {
	module := helper.GetModuleFromRequest(c)
	if module == nil {
		return nil
	}

	systemName := c.Param("system")
	system, err := helper.SystemSvc.GetSystemByName(*module, systemName)
	if err != nil {
		if err == services.ErrSystemNotFound {
			response := ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "System not found",
			}
			c.AbortWithStatusJSON(response.Code, response)
			return nil
		}
		response := ErrorResponse{
			Message: "Unknown error getting system",
		}

		log.Error("error getting system", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return nil
	}

	return system
}

func (helper ApiHelper) GetVersionFromRequest(c *gin.Context) *services.TerraformModuleVersion {
	system := helper.GetSystemFromRequest(c)
	if system == nil {
		return nil
	}

	versionName := c.Param("version")
	version, err := helper.VersionSvc.GetVersionByName(*system, versionName)
	if err != nil {
		if err == services.ErrVersionNotFound {
			response := ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Version not found",
			}
			c.AbortWithStatusJSON(response.Code, response)
			return nil
		}
		response := ErrorResponse{
			Message: "Unknown error getting version",
		}

		log.Error("error getting version", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return nil
	}

	return version
}
