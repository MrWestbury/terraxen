package services

import (
	"context"
	"errors"
	"log"

	"github.com/MrWestbury/terrakube-moduleregistry/backend"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	moduleCollectionName = "module"
)

var (
	ErrModuleNotFound = errors.New("module not found")
)

type TerraformModule struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

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

func (modSvc ModuleService) CreateModule(ns Namespace, module TerraformModule) (*TerraformModule, error) {

	if modSvc.Exists(ns, module.Name) {
		return nil, errors.New("module already exists")
	}

	client := modSvc.Connect()
	ctx := context.Background()
	defer modSvc.HandleDisconnect(client, ctx)

	collection := client.Database(modSvc.Database).Collection(moduleCollectionName)
	_, err := collection.InsertOne(ctx, module)
	if err != nil {
		log.Fatalf("failed to insert module: %v", err)
		return nil, err
	}

	return &module, nil
}

func (modSvc ModuleService) ListModules(ns Namespace) *[]TerraformModule {
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

	return &modules
}

func (modSvc ModuleService) GetModuleByName(ns Namespace, name string) (*TerraformModule, error) {
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

func (modSvc ModuleService) DeleteModule(namespace Namespace, name string) {
	client := modSvc.Connect()
	ctx := context.Background()
	defer modSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"namespace": namespace.Name,
		"name":      name,
	}

	collection := client.Database(modSvc.Database).Collection(moduleCollectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("failed deleting namespace: %v", err)
	}
}

func (modSvc ModuleService) HasChildren(namespace Namespace) bool {
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

func (modSvc ModuleService) Exists(namespace Namespace, module string) bool {
	client := modSvc.Connect()
	ctx := context.Background()
	defer modSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"name":      module,
		"namespace": namespace.Name,
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
