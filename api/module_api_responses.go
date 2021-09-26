package api

import "github.com/MrWestbury/terrakube-moduleregistry/models"

type ResponseListModules struct {
	Modules []models.TerraformModule `json:"modules"`
}
