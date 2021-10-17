package v1

type ResponseVersion struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Download string `json:"downloadUrl"`
}

type ResponseVersionList struct {
	Meta      ListMetaData      `json:"meta"`
	Namespace string            `json:"namespace"`
	Module    string            `json:"module"`
	System    string            `json:"system"`
	Versions  []ResponseVersion `json:"versions"`
}
