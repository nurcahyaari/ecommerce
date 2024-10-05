package database

import (
	"context"

	"github.com/nurcahyaari/ecommerce/config"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	DB *mongo.Database
}

func NewMongoDB(cfg *config.Config) (*MongoDB, error) {
	dbUri := "mongodb://" + cfg.DB.MongoDB.User + ":" + cfg.DB.MongoDB.Pass + "@" + cfg.DB.MongoDB.Host + ":" + string(cfg.DB.MongoDB.Port)
	maxPoolSize := 10
	if cfg.DB.MongoDB.MaxPoolSize != 0 {
		maxPoolSize = cfg.DB.MongoDB.MaxPoolSize
	}

	clientOptions := options.Client()
	clientOptions.ApplyURI(dbUri)
	clientOptions.SetMaxPoolSize(uint64(maxPoolSize))

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, &options.ClientOptions{})
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Info().Str("Name", cfg.DB.MongoDB.Name).Msg("Success connect to MONGODB")
	return &MongoDB{
		DB: client.Database(cfg.DB.MongoDB.Name),
	}, nil
}
