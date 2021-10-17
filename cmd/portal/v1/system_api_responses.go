package v1

import (
	"github.com/MrWestbury/terraxen/services"
)

type RequestNewSystem struct {
	Name string `json:"name"`
}

type ResponseListSystems struct {
	Namespace string                     `json:"namespace"`
	Module    string                     `json:"module"`
	Systems   []services.TerraformSystem `json:"systems"`
}
