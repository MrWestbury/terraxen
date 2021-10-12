package services

import (
	"context"
	"errors"
	"github.com/MrWestbury/terrakube-moduleregistry/backend"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const (
	systemCollectionName = "systems"
)

var (
	ErrSystemNotFound      = errors.New("system not found")
	ErrSystemAlreadyExists = errors.New("system already exists")
)

type TerraformSystem struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Module    string `json:"module"`
}

type SystemService struct {
	backend.MongoBackend
}

func NewSystemService(opts Options) *SystemService {
	svc := &SystemService{
		backend.MongoBackend{
			ConnectionString: opts.ConnectionString(),
			Database:         opts.Database,
		},
	}
	return svc
}

func (sysSvc SystemService) CreateSystem(system TerraformSystem) (*TerraformSystem, error) {
	client := sysSvc.Connect()
	ctx := context.Background()
	defer sysSvc.HandleDisconnect(client, ctx)

	exists := sysSvc.Exists(system)
	if exists {
		return nil, ErrSystemAlreadyExists
	}

	collection := client.Database(sysSvc.Database).Collection(systemCollectionName)
	_, err := collection.InsertOne(ctx, system)
	if err != nil {
		return nil, err
	}

	return &system, nil
}

func (sysSvc SystemService) GetSystemByName(module TerraformModule, systemName string) (*TerraformSystem, error) {
	client := sysSvc.Connect()
	ctx := context.Background()
	defer sysSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"name":      systemName,
		"namespace": module.Namespace,
		"module":    module.Name,
	}

	collection := client.Database(sysSvc.Database).Collection(systemCollectionName)
	var system TerraformSystem
	err := collection.FindOne(ctx, filter).Decode(&system)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSystemNotFound
		}
		log.Printf("unable to check namespace exists: %v", err)
	}
	return &system, nil
}

func (sysSvc SystemService) ListSystemsByModule(module TerraformModule) *[]TerraformSystem {
	client := sysSvc.Connect()
	ctx := context.Background()
	defer sysSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"namespace": module.Namespace,
		"module":    module.Name,
	}

	collection := client.Database(sysSvc.Database).Collection(systemCollectionName)
	var systems []TerraformSystem
	rs, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatalf("unable to check namespace exists: %v", err)
	}

	err = rs.All(ctx, &systems)
	if err != nil {
		log.Fatalf("failed to bind versions objects: %v", err)
	}

	return &systems
}

func (sysSvc SystemService) ExistsByName(module TerraformModule, systemName string) bool {
	client := sysSvc.Connect()
	ctx := context.Background()
	defer sysSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"name":      systemName,
		"namespace": module.Namespace,
		"module":    module.Name,
	}

	collection := client.Database(sysSvc.Database).Collection(moduleCollectionName)
	err := collection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatalf("unable to check namespace exists: %v", err)
	}

	return true
}

func (sysSvc SystemService) Exists(system TerraformSystem) bool {
	client := sysSvc.Connect()
	ctx := context.Background()
	defer sysSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"name":      system.Name,
		"namespace": system.Namespace,
		"module":    system.Name,
	}

	collection := client.Database(sysSvc.Database).Collection(moduleCollectionName)
	err := collection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatalf("unable to check namespace exists: %v", err)
	}

	return true
}
