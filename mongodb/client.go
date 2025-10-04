package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client *mongo.Client

// NewClient creates and returns a new MongoDB client using the provided context and client options.
// It establishes a connection to the MongoDB server.
// Returns the connected client or an error if the connection fails.
func NewClient(ctx context.Context, clientOpts *options.ClientOptions) (Client, error) {
	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	return mongoClient, nil
}
