package models

type TerraformModule struct {
	Versions []ModuleVersion `json:"versions"`
}
