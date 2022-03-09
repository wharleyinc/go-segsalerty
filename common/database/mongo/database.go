package appmongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	URI     string
	Timeout int
}

func NewDriver(config Config) (*mongo.Client, error) {
	if len(config.URI) == 0 {
		return nil, errors.New("invalid_mongo_uri")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		return nil, err
	}

	return client, nil
}
