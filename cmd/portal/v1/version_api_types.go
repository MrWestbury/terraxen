package v1

import "time"

type ResponseVersion struct {
	Id        string    `json:"id"`
	Namespace string    `json:"namespace"`
	Module    string    `json:"module"`
	System    string    `json:"system"`
	Name      string    `json:"name"`
	Download  string    `json:"downloadUrl"`
	Downloads int       `json:"downloads"`
	Created   time.Time `json:"created"`
}

type ResponseVersionList struct {
	Meta     ListMetaData      `json:"meta"`
	Versions []ResponseVersion `json:"versions"`
}
