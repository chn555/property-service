package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds the configuration for MongoDB connection
type Config struct {
	URI     string
	Timeout time.Duration
}

// NewClient creates a new MongoDB client
func NewClient(ctx context.Context, config Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.URI)

	connectCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	client, err := mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}
