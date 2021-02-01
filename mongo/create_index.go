package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func createIndex(db *mongo.Database) (err error) {
	ctx := context.Background()

	_, err = db.Collection(CollState).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{
			Key:   "u",
			Value: 1,
		}},
		Options: options.Index().SetExpireAfterSeconds(int32(time.Hour.Seconds())),
	})
	return
}
