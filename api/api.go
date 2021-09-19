package api

import (
	v1 "github.com/MrWestbury/terrakube-moduleregistry/api/v1"
	"github.com/gin-gonic/gin"
)

type Api struct{}

func (api Api) Router(g *gin.RouterGroup) {
	apirouter := g.Group("api")

	v1api := v1.ApiV1{}
	v1api.Router(apirouter)
}
