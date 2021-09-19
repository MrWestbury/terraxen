package namespace

import (
	"net/http"

	"github.com/MrWestbury/terrakube-moduleregistry/services"
	"github.com/gin-gonic/gin"
)

type NamespaceApi struct {
	service services.NamespaceService
}

func New(svc services.NamespaceService) NamespaceApi {
	new_ns_api := &NamespaceApi{
		service: svc,
	}
	return *new_ns_api
}

func (nsapi NamespaceApi) Router(g *gin.RouterGroup) {
	api_router := g.Group("namespace")

	api_router.GET("/", nsapi.GetNamespaces)                // List namespaces
	api_router.POST("/", nsapi.CreateNamespace)             // New namespace
	api_router.DELETE("/:namespace", nsapi.DeleteNamespace) // Delete namespace
	api_router.GET("/:namespace", nsapi.GetNamespaceByName) // Get specific namespace
}

func (nsapi NamespaceApi) GetNamespaces(c *gin.Context) {
	nslist := nsapi.service.ListNamespaces()

	response := NamespaceListResponse{
		Namespaces: make([]NamespaceResponse, 0),
	}
	for _, ns := range nslist {
		nsr := NamespaceResponse{
			Name:  ns.Name,
			Owner: ns.Owner,
		}
		response.Namespaces = append(response.Namespaces, nsr)
	}

	c.IndentedJSON(http.StatusOK, response)
}

func (nsapi NamespaceApi) CreateNamespace(c *gin.Context) {
	var newNamespace NamespaceResponse

	if err := c.BindJSON(&newNamespace); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	ns, err := nsapi.service.CreateNamespace(newNamespace.Name, newNamespace.Owner)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.IndentedJSON(http.StatusOK, ns)
}

func (nsapi NamespaceApi) GetNamespaceByName(c *gin.Context) {
	namespace := c.Param("namespace")

	ns := nsapi.service.GetNamespaceByName(namespace)
	if ns == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.IndentedJSON(http.StatusOK, ns)
}

func (nsapi NamespaceApi) DeleteNamespace(c *gin.Context) {
	namespace := c.Param("namespace")

	exists := nsapi.service.Exists(namespace)
	if !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	nsapi.service.DeleteNamespace(namespace)
	c.Status(http.StatusAccepted)
}
