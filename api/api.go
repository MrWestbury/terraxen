package api

import (
	"github.com/gin-gonic/gin"
)

type Api struct{}

func (api Api) Router(g *gin.RouterGroup) {
	apirouter := g.Group("api")

	v1api := ApiV1{}
	v1api.Router(apirouter)
}
