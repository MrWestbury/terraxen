package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/MrWestbury/terrakube-moduleregistry/services"
	"github.com/gin-gonic/gin"
)

type NamespaceApi struct {
	service       services.NamespaceService
	moduleService services.ModuleService
}

func NewNamespaceApi(svc services.NamespaceService, modSvc services.ModuleService) *NamespaceApi {
	newNsApi := &NamespaceApi{
		service:       svc,
		moduleService: modSvc,
	}
	return newNsApi
}

func (nsApi NamespaceApi) Router(g *gin.RouterGroup) *gin.RouterGroup {
	apiRouter := g.Group("namespace")

	apiRouter.GET("/", nsApi.GetNamespaces)                // List namespaces
	apiRouter.POST("/", nsApi.CreateNamespace)             // New namespace
	apiRouter.DELETE("/:namespace", nsApi.DeleteNamespace) // Delete namespace
	apiRouter.GET("/:namespace", nsApi.GetNamespaceByName) // Get specific namespace

	return apiRouter
}

func (nsApi NamespaceApi) GetNamespaces(c *gin.Context) {
	nsList := nsApi.service.ListNamespaces()

	response := NamespaceListResponse{
		Namespaces: make([]NamespaceResponse, 0),
	}
	for _, ns := range *nsList {
		nsr := NamespaceResponse{
			Name:  ns.Name,
			Owner: ns.Owner,
		}
		response.Namespaces = append(response.Namespaces, nsr)
	}

	c.IndentedJSON(http.StatusOK, response)
}

func (nsApi NamespaceApi) CreateNamespace(c *gin.Context) {
	var newNamespace NamespaceResponse

	if err := c.BindJSON(&newNamespace); err != nil {
		log.Errorf("create namespace failed: %v", err)
		err := c.AbortWithError(http.StatusNotAcceptable, err)
		if err != nil {
			log.Errorf("error dealing with the error: %v", err)
		}
		return
	}

	ns, err := nsApi.service.CreateNamespace(newNamespace.Name, newNamespace.Owner)
	if err != nil {
		_ = c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.IndentedJSON(http.StatusOK, ns)
}

func (nsApi NamespaceApi) GetNamespaceByName(c *gin.Context) {
	ns := nsApi.GetNamespaceFromRequest(c)
	if ns == nil {
		return
	}

	c.IndentedJSON(http.StatusOK, ns)
}

func (nsApi NamespaceApi) DeleteNamespace(c *gin.Context) {
	namespace := c.Param("namespace")

	ns := nsApi.GetNamespaceFromRequest(c)
	if ns == nil {
		return
	}
	hasChildren := nsApi.moduleService.HasChildren(*ns)
	if hasChildren {
		response := ErrorResponse{
			Code:    http.StatusConflict,
			Message: "namespace has child modules",
		}
		c.AbortWithStatusJSON(response.Code, response)
		return
	}

	nsApi.service.DeleteNamespace(namespace)
	c.Status(http.StatusAccepted)
}

func (nsApi NamespaceApi) GetNamespaceFromRequest(c *gin.Context) *services.Namespace {
	namespaceName := c.Param("namespace")

	exists := nsApi.service.Exists(namespaceName)

	if !exists {
		c.Status(http.StatusNotFound)
		return nil
	}
	ns, err := nsApi.service.GetNamespaceByName(namespaceName)
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
