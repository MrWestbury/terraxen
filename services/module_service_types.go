package services

import (
	"time"
)

type NewTerraformModule struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type TerraformModule struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
