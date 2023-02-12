package storage

import (
	"ChatGo/internal/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	db *mongo.Database
}

func New(ctx context.Context) (*Storage, error) {
	cfg := config.Get()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		return nil, err
	}
	return &Storage{db: client.Database(cfg.Mongo.DB)}, nil
}
