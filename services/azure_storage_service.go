package services

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/MrWestbury/terrakube-moduleregistry/models"
	"io"
	"log"
	"net/url"
)

type AzureStorageService struct {
	accountName   string
	accountKey    string
	containerName string
}

func NewAzureStorageService(accountName string, accountKey string, containerName string) *AzureStorageService {
	svc := &AzureStorageService{
		accountName:   accountName,
		accountKey:    accountKey,
		containerName: containerName,
	}

	return svc
}

func (ass AzureStorageService) UploadStream(version models.ModuleVersion, stream io.ReadSeeker) {
	ctx := context.Background()
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(ass.accountName, ass.accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", ass.accountName, ass.containerName))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)
	// Here's how to upload a blob.
	blobURL := containerURL.NewBlockBlobURL(fmt.Sprintf("%s/%s/%s/terraform-%s-", version.Namespace, version.Module, version.System, version.Version))

	_, err = azblob.UploadStreamToBlockBlob(ctx, stream, blobURL, azblob.UploadStreamToBlockBlobOptions{})
}
