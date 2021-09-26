package models

type ModuleVersion struct {
	Namespace string `json:"namespace"`
	Module    string `json:"module"`
	System    string `json:"System"`
	Version   string `json:"version"`
}
