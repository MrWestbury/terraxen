package v1

import (
	"github.com/MrWestbury/terraxen/services"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NamespaceApi struct {
	helper ApiHelper
}

func NewNamespaceApi(helper ApiHelper) *NamespaceApi {
	newNsApi := &NamespaceApi{
		helper: helper,
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

// GetNamespaces handles gin API request for listing namespaces
func (nsApi NamespaceApi) GetNamespaces(c *gin.Context) {
	nsList := nsApi.helper.NamespaceSvc.ListNamespaces()

	response := ResponseNamespaceList{
		Namespaces: make([]ResponseNamespace, 0),
	}
	for _, ns := range *nsList {
		nsr := ResponseNamespace{
			Name:  ns.Name,
			Owner: ns.Owner,
		}
		response.Namespaces = append(response.Namespaces, nsr)
	}

	c.IndentedJSON(http.StatusOK, response)
}

// CreateNamespace handles gin POST request to create a new namespace
func (nsApi NamespaceApi) CreateNamespace(c *gin.Context) {
	var newNamespace RequestNamespace

	if err := c.BindJSON(&newNamespace); err != nil {
		log.Errorf("create namespace failed: %v", err)
		err := c.AbortWithError(http.StatusNotAcceptable, err)
		if err != nil {
			log.Errorf("error dealing with the error: %v", err)
		}
		return
	}
	newNs := services.NewTerraformNamespace{
		Name:  newNamespace.Name,
		Owner: "Adam",
	}

	ns, err := nsApi.helper.NamespaceSvc.CreateNamespace(newNs)
	if err != nil {
		_ = c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.IndentedJSON(http.StatusOK, ns)
}

func (nsApi NamespaceApi) GetNamespaceByName(c *gin.Context) {
	ns := nsApi.helper.GetNamespaceFromRequest(c)
	if ns == nil {
		return
	}

	c.IndentedJSON(http.StatusOK, ns)
}

func (nsApi NamespaceApi) DeleteNamespace(c *gin.Context) {
	ns := nsApi.helper.GetNamespaceFromRequest(c)
	if ns == nil {
		return
	}

	nsApi.helper.NamespaceSvc.DeleteNamespace(*ns)
	c.Status(http.StatusOK)
}
