package main

import "github.com/MrWestbury/terraxen/services"

type RequestNewModule struct {
	Name string `json:"name"`
}

type ResponseListModules struct {
	Modules []services.TerraformModule `json:"modules"`
}
