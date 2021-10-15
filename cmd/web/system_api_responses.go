package main

import (
	"github.com/MrWestbury/terraxen/services"
)

type ResponseListSystems struct {
	Namespace string                     `json:"namespace"`
	Module    string                     `json:"module"`
	Systems   []services.TerraformSystem `json:"systems"`
}
