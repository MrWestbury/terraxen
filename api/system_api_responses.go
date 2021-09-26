package api

import "github.com/MrWestbury/terrakube-moduleregistry/models"

type ResponseListSystems struct {
	Namespace string                   `json:"namespace"`
	Module    string                   `json:"module"`
	Systems   []models.TerraformSystem `json:"systems"`
}
