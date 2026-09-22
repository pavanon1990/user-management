package mongodb

import (
	"context"
	"fmt"
	"log"
	"user-management-api/pkg/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewClient(ctx context.Context, cfg config.MongoConfig) (*mongo.Client, error) {
	options := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxIdleTime).
		SetConnectTimeout(cfg.ConnectionTimeout)

	client, err := mongo.Connect(options)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectionTimeout)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("%s", fmt.Sprintf("mongo connection error : %s", err))
	}

	log.Printf("connected to mongodb: database=%s", cfg.Database)

	return client, nil
}
