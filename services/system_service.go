package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/MrWestbury/terraxen/backend"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

const (
	systemCollectionName = "systems"
)

var (
	ErrSystemNotFound      = errors.New("system not found")
	ErrSystemAlreadyExists = errors.New("system already exists")
)

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

func (sysSvc SystemService) CreateSystem(systemInfo NewTerraformSystem) (*TerraformSystem, error) {
	newId := fmt.Sprintf("%s/%s/%s", systemInfo.Module.Namespace, systemInfo.Module.Name, systemInfo.Name)
	newSystem := TerraformSystem{
		Id:        newId,
		Name:      systemInfo.Name,
		Namespace: systemInfo.Module.Namespace,
		Module:    systemInfo.Module.Name,
		Created:   time.Now(),
		Updated:   time.Now(),
	}

	client := sysSvc.Connect()
	ctx := context.Background()
	defer sysSvc.HandleDisconnect(client, ctx)

	exists := sysSvc.Exists(newSystem)
	if exists {
		return nil, ErrSystemAlreadyExists
	}

	collection := client.Database(sysSvc.Database).Collection(systemCollectionName)
	_, err := collection.InsertOne(ctx, newSystem)
	if err != nil {
		return nil, err
	}

	return &newSystem, nil
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
		log.Fatalf("unable to check namespace exists: %v", err)
	}

	err = rs.All(ctx, &systems)
	if err != nil {
		log.Fatalf("failed to bind systems objects: %v", err)
	}

	if systems == nil {
		systems = make([]TerraformSystem, 0)
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
		log.Fatalf("unable to check system exists: %v", err)
	}

	return true
}

func (sysSvc SystemService) Exists(system TerraformSystem) bool {
	client := sysSvc.Connect()
	ctx := context.Background()
	defer sysSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"id": system.Id,
	}

	collection := client.Database(sysSvc.Database).Collection(systemCollectionName)
	err := collection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatalf("unable to check system exists: %v", err)
	}

	return true
}

func (sysSvc SystemService) Delete(system TerraformSystem) error {
	client := sysSvc.Connect()
	ctx := context.Background()
	defer sysSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"id": system.Id,
	}

	collection := client.Database(sysSvc.Database).Collection(systemCollectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("failed deleting system: %v", err)
		return err
	}
	return nil
}
