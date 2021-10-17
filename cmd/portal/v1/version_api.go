package v1

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
	helper     ApiHelper
	storageSvc services.AzureStorageService
}

func NewVersionApi(helper ApiHelper, storeSvc services.AzureStorageService) *VersionApi {
	api := &VersionApi{
		helper:     helper,
		storageSvc: storeSvc,
	}

	return api
}

func (versionApi VersionApi) Router(g *gin.RouterGroup) *gin.RouterGroup {
	apiRouter := g.Group("/:system/version")

	apiRouter.GET("/", versionApi.ListVersions)     // List versions
	apiRouter.POST("/", versionApi.CreateHandler)   // New version
	apiRouter.DELETE("/:version", versionApi.Dummy) // Delete version
	apiRouter.GET("/:version", versionApi.Dummy)    // Get specific version

	return apiRouter
}

func (versionApi VersionApi) Dummy(c *gin.Context) {
	response := ErrorResponse{
		Code:    http.StatusNotImplemented,
		Message: "Not yet implemented",
	}

	c.IndentedJSON(http.StatusOK, response)
}

func (versionApi VersionApi) CreateHandler(c *gin.Context) {
	ns := versionApi.helper.GetNamespaceFromRequest(c)
	if ns == nil {
		return
	}

	mod := versionApi.helper.GetModuleFromRequest(c)
	if mod == nil {
		return
	}

	sys := versionApi.helper.GetSystemFromRequest(c)
	if sys == nil {
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

	newVersion := services.NewTerraformModuleVersion{
		Namespace:   ns.Name,
		Module:      mod.Name,
		System:      sys.Name,
		Name:        versionName,
		StoragePath: fmt.Sprintf("%s/%s/%s/%s.zip", ns.Name, mod.Name, sys.Name, versionName),
	}

	// check if the version exists already
	exists := versionApi.helper.VersionSvc.ExistsByName(*sys, versionName)
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
	tmpPath, _ := os.MkdirTemp("", "terraxen")
	defer os.RemoveAll(tmpPath)
	newFileName := fmt.Sprintf("%s/%s%s", tmpPath, uuid.New().String(), extension)
	log.Printf("Save file: %s", newFileName)
	err = c.SaveUploadedFile(uploadFile, newFileName)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fStream, _ := os.Open(newFileName)
	readSeek := io.ReadSeeker(fStream)
	err = versionApi.storageSvc.UploadModuleVersion(newVersion.StoragePath, readSeek)
	if err != nil {
		errBody := ErrorResponse{
			Code:    500,
			Message: "Failed to save version to storage",
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errBody)
		return
	}

	savedVersion, err := versionApi.helper.VersionSvc.CreateVersion(newVersion)
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

// ListVersions handles requests to list available versions of a given system
func (versionApi VersionApi) ListVersions(c *gin.Context) {
	sys := versionApi.helper.GetSystemFromRequest(c)
	if sys == nil {
		return
	}

	versionList, err := versionApi.helper.VersionSvc.ListVersionsBySystem(*sys)
	if err != nil {
		response := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get list of versions",
		}
		c.AbortWithStatusJSON(response.Code, response)
		return
	}

	response := ResponseVersionList{
		Meta: ListMetaData{
			Offset: 0,
			Limit:  0,
		},
		Namespace: sys.Namespace,
		Module:    sys.Module,
		System:    sys.Namespace,
		Versions:  make([]ResponseVersion, 0),
	}
	for _, version := range *versionList {
		response.Versions = append(response.Versions, ResponseVersion{
			Id:       version.Id,
			Name:     version.Name,
			Download: fmt.Sprintf("%s/%s/download", c.Request.URL.String(), version.Name),
		})
	}

	c.IndentedJSON(http.StatusOK, response)
}

func (versionApi VersionApi) DeleteHandler(c *gin.Context) {
	ns := versionApi.helper.GetNamespaceFromRequest(c)
	if ns == nil {
		return
	}

	mod := versionApi.helper.GetModuleFromRequest(c)
	if mod == nil {
		return
	}

	sys := versionApi.helper.GetSystemFromRequest(c)
	if sys == nil {
		return
	}

	ver := versionApi.helper.GetVersionFromRequest(c)
	if ver == nil {
		return
	}

	versionApi.helper.VersionSvc.Delete(*ver)

}
