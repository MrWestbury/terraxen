package services

import (
	"context"
	"errors"
	"log"

	"github.com/MrWestbury/terrakube-moduleregistry/backend"
	"github.com/MrWestbury/terrakube-moduleregistry/models"
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

func (nsSvc NamespaceService) CreateNamespace(name string, owner string) (*models.Namespace, error) {
	if nsSvc.Exists(name) {
		return nil, ErrNamespaceExists
	}

	client := nsSvc.Connect()
	ctx := context.Background()
	defer handleDisconnect(client, ctx)

	newNs := &models.Namespace{
		Name:  name,
		Owner: owner,
	}

	collection := client.Database(nsSvc.Database).Collection(namespaceCollectionName)

	_, err := collection.InsertOne(ctx, newNs)
	if err != nil {
		log.Fatalf("error inserting namespace: %v", err)
	}

	return newNs, nil
}

func (nsSvc NamespaceService) ListNamespaces() *[]models.Namespace {
	client := nsSvc.Connect()
	ctx := context.Background()
	defer handleDisconnect(client, ctx)

	filter := bson.D{}

	collection := client.Database(nsSvc.Database).Collection(namespaceCollectionName)

	var namespaces []models.Namespace
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

func (nsSvc NamespaceService) GetNamespaceByName(name string) (*models.Namespace, error) {
	exists := nsSvc.Exists(name)
	if !exists {
		return nil, errors.New("namespace not found")
	}

	client := nsSvc.Connect()
	ctx := context.Background()
	defer handleDisconnect(client, ctx)

	filter := bson.M{
		"name": name,
	}

	collection := client.Database(nsSvc.Database).Collection(namespaceCollectionName)

	var item models.Namespace
	err := collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNamespaceNotFound
		}
		log.Fatalf("error getting namespace: %v", err)
	}

	return &item, nil
}

func (nsSvc NamespaceService) DeleteNamespace(name string) {
	client := nsSvc.Connect()
	ctx := context.Background()
	defer handleDisconnect(client, ctx)

	filter := bson.M{
		"name": name,
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
	defer handleDisconnect(client, ctx)

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

func handleDisconnect(client *mongo.Client, ctx context.Context) {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatalf("error disconnecting client: %v", err)
	}
}
