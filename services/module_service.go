package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MrWestbury/terraxen/backend"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	moduleCollectionName = "module"
)

var (
	ErrModuleNotFound      = errors.New("module not found")
	ErrModuleAlreadyExists = errors.New("module already exists")
)

type ModuleService struct {
	backend.MongoBackend
}

func NewModuleService(options Options) *ModuleService {
	svc := &ModuleService{
		backend.MongoBackend{
			ConnectionString: options.ConnectionString(),
			Database:         options.Database,
		},
	}

	return svc
}

func (modSvc ModuleService) CreateModule(module NewTerraformModule) (*TerraformModule, error) {
	newId := fmt.Sprintf("%s/%s", module.Namespace, module.Name)
	newModule := TerraformModule{
		Id:        newId,
		Name:      module.Name,
		Namespace: module.Namespace,
		Created:   time.Now(),
		Updated:   time.Now(),
	}

	if modSvc.Exists(newModule) {
		return nil, ErrModuleAlreadyExists
	}

	client := modSvc.Connect()
	ctx := context.Background()
	defer modSvc.HandleDisconnect(client, ctx)

	collection := client.Database(modSvc.Database).Collection(moduleCollectionName)
	_, err := collection.InsertOne(ctx, newModule)
	if err != nil {
		log.Fatalf("failed to insert module: %v", err)
		return nil, err
	}

	return &newModule, nil
}

func (modSvc ModuleService) ListModules(ns TerraformNamespace) *[]TerraformModule {
	client := modSvc.Connect()
	ctx := context.Background()
	defer modSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"namespace": ns.Name,
	}

	collection := client.Database(modSvc.Database).Collection(moduleCollectionName)

	var modules []TerraformModule
	rs, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("failed getting modules: %v", err)
	}

	err = rs.All(ctx, &modules)
	if err != nil {
		log.Fatalf("failed decoding modules: %v", err)
	}

	if modules == nil {
		modules = make([]TerraformModule, 0)
	}

	return &modules
}

func (modSvc ModuleService) GetModuleByName(ns TerraformNamespace, name string) (*TerraformModule, error) {
	client := modSvc.Connect()
	ctx := context.Background()
	defer modSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"namespace": ns.Name,
		"name":      name,
	}

	collection := client.Database(modSvc.Database).Collection(moduleCollectionName)

	var item TerraformModule
	err := collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrModuleNotFound
		}
		log.Fatalf("error getting module: %v", err)
	}

	return &item, nil
}

func (modSvc ModuleService) DeleteModule(module TerraformModule) {
	client := modSvc.Connect()
	ctx := context.Background()
	defer modSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"_id": module.Id,
	}

	collection := client.Database(modSvc.Database).Collection(moduleCollectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("failed deleting module: %v", err)
	}
}

func (modSvc ModuleService) HasChildren(namespace TerraformNamespace) bool {
	client := modSvc.Connect()
	ctx := context.Background()
	defer modSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"namespace": namespace.Name,
	}

	collection := client.Database(modSvc.Database).Collection(moduleCollectionName)
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatalf("error getting module count: %v", err)
	}
	return count > 0
}

func (modSvc ModuleService) Exists(module TerraformModule) bool {
	client := modSvc.Connect()
	ctx := context.Background()
	defer modSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"id": module.Id,
	}

	collection := client.Database(modSvc.Database).Collection(moduleCollectionName)
	err := collection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatalf("unable to check namespace exists: %v", err)
	}

	return true
}
