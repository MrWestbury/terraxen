package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/MrWestbury/terraxen/backend"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

const (
	versionCollectionName = "versions"
)

var (
	ErrVersionNotFound      = errors.New("version not found for module")
	ErrVersionAlreadyExists = errors.New("version already exists")
)

type VersionService struct {
	backend.MongoBackend
}

func NewVersionService(opts Options) *VersionService {
	svc := &VersionService{
		backend.MongoBackend{
			ConnectionString: opts.ConnectionString(),
			Database:         opts.Database,
		},
	}
	return svc
}

func (verSvc VersionService) ListVersionsBySystem(system TerraformSystem) (*[]TerraformModuleVersion, error) {
	client := verSvc.Connect()
	ctx := context.Background()
	defer verSvc.HandleDisconnect(client, ctx)

	var results []TerraformModuleVersion

	filter := bson.M{
		"namespace": system.Namespace,
		"module":    system.Module,
		"system":    system.Name,
	}

	collection := client.Database(verSvc.Database).Collection(versionCollectionName)

	rs, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("failed getting list of versions: %v", err)
		return nil, err
	}

	err = rs.All(ctx, &results)
	if err != nil {
		log.Fatalf("failed decoding versions: %v", err)
		return nil, err
	}

	return &results, nil
}

func (verSvc VersionService) Exists(version TerraformModuleVersion) bool {
	client := verSvc.Connect()
	ctx := context.Background()
	defer verSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"namespace": version.Namespace,
		"module":    version.Module,
		"system":    version.System,
		"name":      version.Name,
	}

	collection := client.Database(verSvc.Database).Collection(versionCollectionName)

	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return false
	}
	return true
}

func (verSvc VersionService) ExistsByName(system TerraformSystem, versionName string) bool {
	client := verSvc.Connect()
	ctx := context.Background()
	defer verSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"namespace": system.Namespace,
		"module":    system.Module,
		"system":    system.Name,
		"name":      versionName,
	}

	collection := client.Database(verSvc.Database).Collection(versionCollectionName)

	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return false
	}
	return true
}

func (verSvc VersionService) GetVersionByName(system TerraformSystem, versionName string) (*TerraformModuleVersion, error) {
	client := verSvc.Connect()
	ctx := context.Background()
	defer verSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"namespace": system.Namespace,
		"module":    system.Module,
		"system":    system.Name,
		"name":      versionName,
	}

	collection := client.Database(verSvc.Database).Collection(versionCollectionName)

	var version TerraformModuleVersion
	err := collection.FindOne(ctx, filter).Decode(&version)
	if err != nil {
		return nil, ErrVersionNotFound
	}
	return &version, nil
}

func (verSvc VersionService) CreateVersion(newVersion NewTerraformModuleVersion) (*TerraformModuleVersion, error) {
	id := fmt.Sprintf("%s/%s/%s/%s", newVersion.Namespace, newVersion.Module, newVersion.System, newVersion.Name)
	version := &TerraformModuleVersion{
		Id:          id,
		Namespace:   newVersion.Namespace,
		Module:      newVersion.Module,
		System:      newVersion.System,
		Name:        newVersion.Name,
		Created:     time.Now(),
		DownloadKey: newVersion.StoragePath,
	}

	c := verSvc.Connect()
	ctx := context.Background()
	defer verSvc.HandleDisconnect(c, ctx)

	exists := verSvc.Exists(*version)
	if exists {
		return nil, ErrVersionAlreadyExists
	}

	nsCollection := c.Database(verSvc.Database).Collection(versionCollectionName)
	_, err := nsCollection.InsertOne(ctx, version)
	if err != nil {
		log.Printf("failed to add version: %v", err)
		return nil, err
	}

	return version, nil
}

func (verSvc VersionService) Delete(version TerraformModuleVersion) {
	client := verSvc.Connect()
	ctx := context.Background()
	defer verSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"id": version.Id,
	}

	collection := client.Database(verSvc.Database).Collection(versionCollectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("failed deleting namespace: %v", err)
	}
}

func (verSvc VersionService) HasChildren(system TerraformSystem) bool {
	client := verSvc.Connect()
	ctx := context.Background()
	defer verSvc.HandleDisconnect(client, ctx)

	filter := bson.M{
		"namespace": system.Namespace,
		"module":    system.Module,
		"system":    system.Name,
	}

	collection := client.Database(verSvc.Database).Collection(moduleCollectionName)
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatalf("error getting system count: %v", err)
	}
	return count > 0
}

func (verSvc VersionService) GetDownloadCount(version TerraformModuleVersion) int {
	return 100
}
