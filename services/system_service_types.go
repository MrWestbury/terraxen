package services

import (
	"time"
)

type NewTerraformSystem struct {
	Name   string
	Module TerraformModule
}

type TerraformSystem struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Module    string    `json:"module"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
