package main

import (
	"fmt"
	"github.com/MrWestbury/terraxen/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ModuleRegistryRootPath = "/module/v1"
)

type ResponseVersion struct {
	Version string `json:"version"`
}

type ResponseModule struct {
	Versions []ResponseVersion `json:"versions"`
}

type ResponseModuleParent struct {
	Modules []ResponseModule `json:"modules"`
}

type ModuleRegistry struct {
	nsSvc      services.NamespaceService
	modSvc     services.ModuleService
	sysSvc     services.SystemService
	verSvc     services.VersionService
	storageSvc services.AzureStorageService
}

func NewModuleRegistry(nssvc services.NamespaceService, modsvc services.ModuleService, syssvc services.SystemService, versvc services.VersionService, storageSvc services.AzureStorageService) *ModuleRegistry {
	reg := &ModuleRegistry{
		nsSvc:      nssvc,
		modSvc:     modsvc,
		sysSvc:     syssvc,
		verSvc:     versvc,
		storageSvc: storageSvc,
	}

	return reg
}

func (reg ModuleRegistry) Router(g *gin.RouterGroup) *gin.RouterGroup {
	regRouter := g.Group(ModuleRegistryRootPath)

	regRouter.GET(":namespace/:name/:system/versions", reg.getModuleVersions)
	regRouter.GET(":namespace/:name/:system/:version/download", reg.getVersionsDownload)
	regRouter.GET(":namespace/:name/:system/:version/downloadFile")
	return regRouter
}

func (reg ModuleRegistry) getModuleVersions(c *gin.Context) {
	namespace := c.Param("namespace")
	ns, err := reg.nsSvc.GetNamespaceByName(namespace)
	if err == services.ErrNamespaceNotFound {
		body := ErrorResponse{
			Message: "Namespace not found",
		}
		c.AbortWithStatusJSON(http.StatusNotFound, body)
	}

	// Namespace exists, now check module
	modName := c.Param("module")
	module, err := reg.modSvc.GetModuleByName(*ns, modName)
	if err != nil {
		if err == services.ErrModuleNotFound {
			response := ErrorResponse{
				Message: "Module not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, response)
		}
		response := ErrorResponse{
			Message: "unknown error getting module",
		}
		c.AbortWithStatusJSON(http.StatusNotFound, response)
	}

	// Module exists, check system
	sysName := c.Param("system")
	system, err := reg.sysSvc.GetSystemByName(*module, sysName)
	if err != nil {
		if err == services.ErrSystemNotFound {
			body := ErrorResponse{
				Message: "System not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
		body := ErrorResponse{
			Message: "Unknown error getting system",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, body)
	}

	// System exists, finally, get the versions
	versions := reg.verSvc.ListVersionsBySystem(*system)
	if err != nil {

	}

	var versionsList []ResponseVersion
	for _, ver := range *versions {
		newVer := ResponseVersion{
			Version: ver.Name,
		}
		versionsList = append(versionsList, newVer)
	}
	modResponse := ResponseModule{
		Versions: versionsList,
	}
	response := ResponseModuleParent{
		Modules: []ResponseModule{modResponse},
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (reg ModuleRegistry) getVersionsDownload(c *gin.Context) {
	namespace := c.Param("namespace")
	ns, err := reg.nsSvc.GetNamespaceByName(namespace)
	if err == services.ErrNamespaceNotFound {
		body := map[string]string{
			"error": "Namespace not found",
		}
		c.IndentedJSON(http.StatusNotFound, body)
		return
	}

	// Namespace exists, now check module
	modName := c.Param("module")
	module, err := reg.modSvc.GetModuleByName(*ns, modName)
	if err != nil {
		if err == services.ErrModuleNotFound {
			response := ErrorResponse{
				Message: "Module not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, response)
		}
		response := ErrorResponse{
			Message: "Unknown error get module details",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
	}

	// Module exists, check system
	sysName := c.Param("system")
	system, err := reg.sysSvc.GetSystemByName(*module, sysName)
	if err != nil {
		if err == services.ErrSystemNotFound {
			body := ErrorResponse{
				Message: "System not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
		body := ErrorResponse{
			Message: "Unknown error getting system",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, body)
	}

	versionName := c.Param("version")
	version, err := reg.verSvc.GetVersionByName(*system, versionName)
	if err != nil {
		if err == services.ErrVersionNotFound {
			body := ErrorResponse{
				Message: "Module version not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, body)
		}
	}
	dlPath := fmt.Sprintf("%s/%s/%s/%s/%s/downloadFile", ModuleRegistryRootPath, version.Namespace, version.Module, version.System, version.Name)

	c.Header("X-Terraform-Get", dlPath)
	c.Status(http.StatusOK)
}
