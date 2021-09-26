package models

type TerraformSystem struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Module    string `json:"module"`
}
