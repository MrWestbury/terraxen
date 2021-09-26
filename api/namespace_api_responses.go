package api

type NamespaceResponse struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type NamespaceListResponse struct {
	Namespaces []NamespaceResponse `json:"namespaces"`
}
