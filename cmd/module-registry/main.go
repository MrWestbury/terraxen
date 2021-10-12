package main

import (
	"github.com/MrWestbury/terrakube-moduleregistry/services"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	router := gin.Default()
	rootGroup := router.Group("/")

	azureOpts := services.AzureStorageOptions{
		AccountName:   os.Getenv("AZURE_STORAGE_ACCOUNT"),
		AccountKey:    os.Getenv("AZURE_STORAGE_KEY"),
		ContainerName: os.Getenv("AZURE_STORAGE_CONTAINER_NAME"),
	}

	storageSvc := services.NewAzureStorageService(azureOpts)

	svcOpts := services.Options{
		Hostname: os.Getenv("MONGO_HOSTNAME"),
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Database: os.Getenv("MONGO_DB"),
		ExtraOpts: map[string]string{
			"retryWrites": "true",
			"w":           "majority",
		},
	}

	nsSvc := services.NewNamespaceService(svcOpts)
	modSvc := services.NewModuleService(svcOpts)
	sysSvc := services.NewSystemService(svcOpts)
	verSvc := services.NewVersionService(svcOpts)

	moduleRegistry := NewModuleRegistry(*nsSvc, *modSvc, *sysSvc, *verSvc, *storageSvc)
	moduleRegistry.Router(rootGroup)

	err := router.Run(":8010")
	if err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
