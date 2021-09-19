package v1

import (
	"github.com/MrWestbury/terrakube-moduleregistry/api/v1/namespace"
	"github.com/MrWestbury/terrakube-moduleregistry/services"
	"github.com/gin-gonic/gin"
)

type ApiV1 struct{}

func (v1 ApiV1) Router(g *gin.RouterGroup) {
	group := g.Group("v1")

	ns_svc := services.NamespaceService{}.New()

	ns_api := namespace.New(ns_svc)
	ns_api.Router(group)
}
