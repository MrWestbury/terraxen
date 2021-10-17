package v1

import (
	"github.com/MrWestbury/terraxen/services"
	"github.com/gin-gonic/gin"
	"os"
)

type V1 struct{}

func (v1 V1) Router(g *gin.RouterGroup) {
	group := g.Group("v1")

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

	helper := NewHelper(*nsSvc, *modSvc, *sysSvc, *verSvc)

	nsApi := NewNamespaceApi(*helper)
	modApi := NewModuleApi(*helper)
	sysApi := NewSystemApi(*helper)
	verApi := NewVersionApi(*helper, *storageSvc)

	nsGroup := nsApi.Router(group)
	modGroup := modApi.Router(nsGroup)
	sysGroup := sysApi.Router(modGroup)
	verApi.Router(sysGroup)

}
