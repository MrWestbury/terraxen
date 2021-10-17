package main

import (
	"github.com/MrWestbury/terraxen/cmd/portal/v1"
	"github.com/gin-gonic/gin"
)

type Api struct{}

func NewApi() *Api {
	var newApi = &Api{}

	return newApi
}

func (api Api) Router(g *gin.RouterGroup) {
	apiRouter := g.Group("/api")

	v1api := v1.V1{}
	v1api.Router(apiRouter)
}
