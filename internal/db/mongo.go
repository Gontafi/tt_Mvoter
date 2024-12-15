package db

import (
	"context"
	"tt/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongoDB(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	if cfg.MongoUsername != "" && cfg.MongoPassword != "" {
		clientOptions.SetAuth(options.Credential{
			Username: cfg.MongoUsername,
			Password: cfg.MongoPassword,
		})
	}

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
