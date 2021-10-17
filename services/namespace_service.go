package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/MrWestbury/terraxen/backend"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	namespaceCollectionName = "namespace"
)

var (
	ErrNamespaceNotFound = errors.New("namespace not found")
	ErrNamespaceExists   = errors.New("namespace already exists")
)

type NamespaceService struct {
	backend.MongoBackend
}

func NewNamespaceService(options Options) *NamespaceService {
	svc := &NamespaceService{
		backend.MongoBackend{
			ConnectionString: options.ConnectionString(),
			Database:         options.Database,
		},
	}

	return svc
}

func (nsSvc NamespaceService) CreateNamespace(newNs NewTerraformNamespace) (*TerraformNamespace, error) {
	if nsSvc.Exists(newNs.Name) {
		return nil, ErrNamespaceExists
	}

	client := nsSvc.Connect()
	ctx := context.Background()
	defer nsSvc.HandleDisconnect(client, ctx)

	newId := newNs.Name
	ns := &TerraformNamespace{
		Id:      newId,
		Name:    newNs.Name,
		Owner:   newNs.Owner,
		Created: time.Now(),
		Updated: time.Now(),
	}

	collection := client.Database(nsSvc.Database).Collection(namespaceCollectionName)

	_, err := collection.InsertOne(ctx, ns)
	if err != nil {
		log.Fatalf("error inserting namespace: %v", err)
	}

	return ns, nil
}

func (nsSvc NamespaceService) ListNamespaces() *[]TerraformNamespace {
	client := nsSvc.Connect()
	ctx := context.Background()
	defer nsSvc.HandleDisconnect(client, ctx)

	filter := bson.D{}

	collection := client.Database(nsSvc.Database).Collection(namespaceCollectionName)

	var namespaces []TerraformNamespace
	rs, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("failed getting namespaces: %v", err)
	}

	err = rs.All(ctx, &namespaces)
	if err != nil {
		log.Fatalf("failed decoding namespaces: %v", err)
	}

	return &namespaces
}

func (nsSvc NamespaceService) GetNamespaceByName(name string) (*TerraformNamespace, error) {
	exists := nsSvc.Exists(name)
	if !exists {
		return nil, errors.New("namespace not found")
	}

	client := nsSvc.Connect()
	ctx := context.Background()
	defer nsSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"name": name,
	}

	collection := client.Database(nsSvc.Database).Collection(namespaceCollectionName)

	var item TerraformNamespace
	err := collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNamespaceNotFound
		}
		log.Fatalf("error getting namespace: %v", err)
	}

	return &item, nil
}

func (nsSvc NamespaceService) DeleteNamespace(namespace TerraformNamespace) {
	client := nsSvc.Connect()
	ctx := context.Background()
	defer nsSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"id": namespace.Id,
	}

	collection := client.Database(nsSvc.Database).Collection(namespaceCollectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("failed deleting namespace: %v", err)
	}
}

func (nsSvc NamespaceService) Exists(name string) bool {
	client := nsSvc.Connect()
	ctx := context.Background()
	defer nsSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"name": name,
	}

	collection := client.Database(nsSvc.Database).Collection(namespaceCollectionName)
	err := collection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatalf("unable to check namespace exists: %v", err)
	}

	return true
}
