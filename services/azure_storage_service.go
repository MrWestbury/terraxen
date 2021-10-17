package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"
)

type AzureStorageOptions struct {
	AccountName   string
	AccountKey    string
	ContainerName string
}

type AzureStorageService struct {
	accountName   string
	accountKey    string
	containerName string
}

func NewAzureStorageService(options AzureStorageOptions) *AzureStorageService {
	svc := &AzureStorageService{
		accountName:   options.AccountName,
		accountKey:    options.AccountKey,
		containerName: options.ContainerName,
	}

	return svc
}

func (ass AzureStorageService) DownloadModuleVersion(c *gin.Context, version TerraformModuleVersion) {
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(ass.accountName, ass.accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", ass.accountName, ass.containerName))
	containerURL := azblob.NewContainerURL(*URL, p)

	blobURL := containerURL.NewBlockBlobURL(version.DownloadKey)

	ctx := context.Background()
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})

	if err != nil {
		log.Printf("Failed to get download stream: %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	data := &bytes.Buffer{}
	_, _ = data.ReadFrom(bodyStream)

	c.Stream(func(w io.Writer) bool {
		_, _ = w.Write(data.Next(data.Len()))
		return false
	})
}

func (ass AzureStorageService) UploadModuleVersion(path string, stream io.ReadSeeker) error {
	// Create a request pipeline using your Storage account's name and account key.
	credential, err := azblob.NewSharedKeyCredential(ass.accountName, ass.accountKey)
	if err != nil {
		log.Fatal(err)
		return err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your Storage account blob service URL endpoint.
	cURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", ass.accountName, ass.containerName))

	// Create an ServiceURL object that wraps the service URL and a request pipeline to making requests.
	containerURL := azblob.NewContainerURL(*cURL, p)

	ctx := context.Background() // This example uses a never-expiring context
	// Here's how to create a blob with HTTP headers and metadata (I'm using the same metadata that was put on the container):

	blobURL := containerURL.NewBlockBlobURL(path)

	// Wrap the request body in a RequestBodyProgress and pass a callback function for progress reporting.
	_, err = blobURL.Upload(ctx, stream, azblob.BlobHTTPHeaders{
		ContentType:        "text/html; charset=utf-8",
		ContentDisposition: "attachment",
	}, azblob.Metadata{}, azblob.BlobAccessConditions{}, azblob.DefaultAccessTier, nil, azblob.ClientProvidedKeyOptions{})

	if err != nil {
		log.Fatalf("failed uploading file: %v", err)
		return err
	}

	return nil
}
