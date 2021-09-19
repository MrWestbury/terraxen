package services

import "errors"

type NamespaceService struct {
	tmpstore map[string]Namespace
}

type Namespace struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func (nssvc NamespaceService) New() NamespaceService {
	ns_svc := NamespaceService{}
	ns_svc.tmpstore = make(map[string]Namespace)
	return ns_svc
}

func (nssvc NamespaceService) CreateNamespace(name string, owner string) (*Namespace, error) {
	if nssvc.Exists(name) {
		return nil, errors.New("Namespace already exists")
	}

	new_ns := &Namespace{
		Name:  name,
		Owner: owner,
	}

	nssvc.tmpstore[name] = *new_ns

	return new_ns, nil
}

func (nssvc NamespaceService) ListNamespaces() []Namespace {
	var results []Namespace
	for _, obj := range nssvc.tmpstore {
		results = append(results, obj)
	}
	return results
}

func (nssvc NamespaceService) GetNamespaceByName(name string) *Namespace {
	if !nssvc.Exists(name) {
		return nil
	}

	item := nssvc.tmpstore[name]
	return &item
}

func (nssvc NamespaceService) DeleteNamespace(name string) {
	delete(nssvc.tmpstore, name)
}

func (nssvc NamespaceService) Exists(name string) bool {
	_, found := nssvc.tmpstore[name]
	return found
}
