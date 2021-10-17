package services

import "time"

type NewTerraformModuleVersion struct {
	Namespace   string `json:"namespace"`
	Module      string `json:"module"`
	System      string `json:"system"`
	Name        string `json:"name"`
	StoragePath string
}

type TerraformModuleVersion struct {
	Id          string    `json:"_id"`
	Namespace   string    `json:"namespace"`
	Module      string    `json:"module"`
	System      string    `json:"system"`
	Name        string    `json:"name"`
	DownloadKey string    `json:"downloadkey"`
	Created     time.Time `json:"created"`
}
