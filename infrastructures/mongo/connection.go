package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDB(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second) // much shorter response time
	defer pingCancel()

	err = client.Ping(pingCtx, readpref.Primary()) // Check primary member of replica set in the cluster
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB primary after connection: %w", err)
	}
	return client, nil
}
