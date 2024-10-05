package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CartRepository struct {
	db             *database.MongoDB
	collectionName string
}

func NewCartRepository(db *database.MongoDB) repository.CartRepositorier {
	return &CartRepository{
		db:             db,
		collectionName: "cart",
	}
}

func (r *CartRepository) FindCart(ctx context.Context, filter entity.CartFilter) (entity.Carts, error) {
	collection := r.db.DB.Collection(r.collectionName)
	cursor, err := collection.Find(ctx, filter.Filter())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	carts := make(entity.Carts, 0)
	for cursor.Next(ctx) {
		var cart entity.Cart
		err := cursor.Decode(&cart)
		if err != nil {
			return nil, err
		}

		carts = append(carts, cart)
	}

	return carts, nil
}

func (r *CartRepository) UpsertCart(ctx context.Context, cart entity.Cart) error {
	collection := r.db.DB.Collection(r.collectionName)
	err := collection.FindOneAndUpdate(
		ctx,
		bson.M{
			"_id": cart.Id,
		},
		bson.M{"$set": cart},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Err()
	if err != nil && err == mongo.ErrNoDocuments {
		collection.InsertOne(ctx, cart)
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
