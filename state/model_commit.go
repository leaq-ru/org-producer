package state

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (m Model) Commit(ctx context.Context, lastLoopIndex uint32) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = m.coll.UpdateOne(ctx, bson.M{}, bson.M{
		"$set": state{
			LastLoopIndex: lastLoopIndex,
			UpdatedAt:     time.Now().UTC(),
		},
	}, options.Update().SetUpsert(true))
	return
}
