package services

import (
	"context"
	"errors"
	"github.com/MrWestbury/terraxen/backend"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

const (
	versionCollectionName = "versions"
)

var (
	ErrVersionNotFound      = errors.New("version not found for module")
	ErrVersionAlreadyExists = errors.New("version already exists")
)

type ModuleVersion struct {
	Namespace   string `json:"namespace"`
	Module      string `json:"module"`
	System      string `json:"system"`
	Name        string `json:"name"`
	DownloadKey string `json:"downloadkey"`
}

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

func (verSvc VersionService) ListVersionsBySystem(system TerraformSystem) *[]ModuleVersion {
	client := verSvc.Connect()
	ctx := context.Background()
	defer verSvc.HandleDisconnect(client, ctx)

	var results []ModuleVersion

	filter := bson.M{
		"namespace": system.Namespace,
		"module":    system.Module,
		"system":    system.Name,
	}

	collection := client.Database(verSvc.Database).Collection(versionCollectionName)

	rs, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("failed getting list of versions: %v", err)
	}

	err = rs.All(ctx, &results)
	if err != nil {
		log.Fatalf("failed decoding versions: %v", err)
	}

	return &results
}

func (verSvc VersionService) Exists(version ModuleVersion) bool {
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

func (verSvc VersionService) GetVersionByName(system TerraformSystem, versionName string) (*ModuleVersion, error) {
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

	var version ModuleVersion
	err := collection.FindOne(ctx, filter).Decode(&version)
	if err != nil {
		return nil, ErrVersionNotFound
	}
	return &version, nil
}

func (verSvc VersionService) CreateVersion(version ModuleVersion) (*ModuleVersion, error) {
	c := verSvc.Connect()
	ctx := context.Background()
	defer verSvc.HandleDisconnect(c, ctx)

	exists := verSvc.Exists(version)
	if exists {
		return nil, ErrVersionAlreadyExists
	}

	nsCollection := c.Database(verSvc.Database).Collection(versionCollectionName)
	_, err := nsCollection.InsertOne(ctx, version)
	if err != nil {
		log.Printf("failed to add version: %v", err)
		return nil, err
	}

	return &version, nil
}
