package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MrWestbury/terraxen/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VersionApi struct {
	namespaceSvc services.NamespaceService
	moduleSvc    services.ModuleService
	systemSvc    services.SystemService
	versionSvc   services.VersionService
	storageSvc   services.AzureStorageService
}

func NewVersionApi(nssvc services.NamespaceService, modsvc services.ModuleService, syssvc services.SystemService, versvc services.VersionService, storeSvc services.AzureStorageService) *VersionApi {
	api := &VersionApi{
		namespaceSvc: nssvc,
		moduleSvc:    modsvc,
		systemSvc:    syssvc,
		versionSvc:   versvc,
		storageSvc:   storeSvc,
	}

	return api
}

func (versionApi VersionApi) Router(g *gin.RouterGroup) *gin.RouterGroup {
	apiRouter := g.Group("/:system/version")

	apiRouter.GET("/", versionApi.Dummy)            // List namespaces
	apiRouter.POST("/", versionApi.Dummy)           // New namespace
	apiRouter.DELETE("/:version", versionApi.Dummy) // Delete namespace
	apiRouter.GET("/:version", versionApi.Dummy)    // Get specific namespace

	return apiRouter
}

func (versionApi VersionApi) Dummy(c *gin.Context) {
	response := ErrorResponse{
		Code:    http.StatusNotImplemented,
		Message: "Not yet implemented",
	}

	c.IndentedJSON(http.StatusOK, response)
}

func (versionApi VersionApi) ListHandler(c *gin.Context) {
	response := ErrorResponse{
		Code:    http.StatusNotImplemented,
		Message: "Not yet implemented",
	}

	c.IndentedJSON(http.StatusOK, response)
}

func (versionApi VersionApi) CreateHandler(c *gin.Context) {
	namespaceName := c.Param("namespace")
	moduleName := c.Param("module")
	systemName := c.Param("system")

	// Check namespace
	ns, err := versionApi.namespaceSvc.GetNamespaceByName(namespaceName)
	if err != nil {
		if err == services.ErrNamespaceNotFound {
			errBody := ErrorResponse{
				Code:    404,
				Message: "Namespace not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, errBody)
			return
		}
		errBody := ErrorResponse{
			Code:    500,
			Message: "Namespace verification failed due to internal error",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errBody)
		return
	}

	// Check module
	mod, err := versionApi.moduleSvc.GetModuleByName(*ns, moduleName)
	if err != nil {
		if err == services.ErrModuleNotFound {
			errBody := ErrorResponse{
				Code:    404,
				Message: "Module not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, errBody)
			return
		}
		errBody := ErrorResponse{
			Code:    500,
			Message: "Module verification failed due to internal error",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errBody)
		return
	}

	// Check system
	sys, err := versionApi.systemSvc.GetSystemByName(*mod, systemName)
	if err != nil {
		if err == services.ErrSystemNotFound {
			errBody := ErrorResponse{
				Code:    404,
				Message: "System not found",
			}
			c.AbortWithStatusJSON(http.StatusNotFound, errBody)
			return
		}
		errBody := ErrorResponse{
			Code:    500,
			Message: "System verification failed due to internal error",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errBody)
		return
	}

	uploadFile, err := c.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		errBody := ErrorResponse{
			Code:    500,
			Message: "Failed to upload file due to internal error",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errBody)
		return
	}

	versionName := c.PostForm("version")

	newVersion := services.ModuleVersion{
		Namespace: ns.Name,
		Module:    mod.Name,
		System:    sys.Name,
		Name:      versionName,
	}

	// check if the version exists already
	exists := versionApi.versionSvc.Exists(newVersion)
	if err != nil {
		log.Printf("Failed to get check version exists: %v", err)
		errBody := ErrorResponse{
			Code:    500,
			Message: "Failed to upload file due to internal error",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errBody)
		return
	}

	if exists {

		errBody := ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Version already exists",
		}
		c.AbortWithStatusJSON(http.StatusConflict, errBody)
		return
	}

	extension := filepath.Ext(uploadFile.Filename)
	newFileName := fmt.Sprintf("/tmp/%s%s", uuid.New().String(), extension)
	err = c.SaveUploadedFile(uploadFile, newFileName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	fStream, _ := os.Open(newFileName)
	readseek := io.ReadSeeker(fStream)
	version := versionApi.storageSvc.UploadModuleVersion(newVersion, readseek)

	savedVersion, err := versionApi.versionSvc.CreateVersion(version)
	if err != nil {
		errBody := ErrorResponse{
			Code:    500,
			Message: "Failed to save version",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errBody)
		return
	}

	c.IndentedJSON(http.StatusOK, savedVersion)
}
