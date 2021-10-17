package v1

type RequestNamespace struct {
	Name string `json:"name"`
}

type ResponseNamespace struct {
	Name  string   `json:"name"`
	Owner string   `json:"owner"`
	Links []string `json:"links"`
}

type ResponseNamespaceList struct {
	Namespaces []ResponseNamespace `json:"namespaces"`
}
