package main

import "github.com/MrWestbury/terrakube-moduleregistry/services"

type ResponseListModules struct {
	Modules []services.TerraformModule `json:"modules"`
}
