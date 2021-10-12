package backend

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBackend struct {
	ConnectionString string
	Database         string
}

func (mb MongoBackend) Connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mb.ConnectionString) //.SetDirect(true)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("unable to initialize connection %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("unable to connect %v", err)
	}

	return client
}

func (mb MongoBackend) HandleDisconnect(client *mongo.Client, ctx context.Context) {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatalf("error disconnecting client: %v", err)
	}
}
