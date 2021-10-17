package services

import (
	"time"
)

type NewTerraformNamespace struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type TerraformNamespace struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Owner   string    `json:"owner"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
