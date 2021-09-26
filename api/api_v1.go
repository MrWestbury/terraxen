package api

import (
	"github.com/MrWestbury/terrakube-moduleregistry/services"
	"github.com/gin-gonic/gin"
	"os"
)

type ApiV1 struct{}

func (v1 ApiV1) Router(g *gin.RouterGroup) {
	group := g.Group("v1")

	svcOpts := services.Options{
		Hostname: os.Getenv("MONGO_HOSTNAME"),
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Database: os.Getenv("MONGO_DB"),
		ExtraOpts: map[string]string{
			"retryWrites": "true",
			"w":           "majority",
		},
	}

	nsSvc := services.NewNamespaceService(svcOpts)
	modSvc := services.NewModuleService(svcOpts)
	sysSvc := services.NewSystemService(svcOpts)
	verSvc := services.NewVersionService(svcOpts)

	nsApi := NewNamespaceApi(*nsSvc)
	modApi := NewModuleApi(*nsSvc, *modSvc)
	sysApi := NewSystemApi(*nsSvc, *modSvc, *sysSvc)
	verApi := NewVersionApi(*nsSvc, *modSvc, *sysSvc, *verSvc)

	ns_group := nsApi.Router(group)
	mod_group := modApi.Router(ns_group)
	sys_group := sysApi.Router(mod_group)
	verApi.Router(sys_group)

}
