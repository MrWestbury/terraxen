package v1

import "github.com/MrWestbury/terraxen/services"

type RequestNewModule struct {
	Name string `json:"name"`
}

type ResponseListModules struct {
	Modules []services.TerraformModule `json:"modules"`
}

type ResponseTerraformModule struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
