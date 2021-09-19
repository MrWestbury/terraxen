package models

type ModuleVersionsParent struct {
	Modules []TerraformModule `json:"modules"`
}
