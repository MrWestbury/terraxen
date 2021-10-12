package main

import (
	"github.com/MrWestbury/terrakube-moduleregistry/services"
)

type ResponseListSystems struct {
	Namespace string                     `json:"namespace"`
	Module    string                     `json:"module"`
	Systems   []services.TerraformSystem `json:"systems"`
}
