package state

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (m Model) GetLastLoopIndex(ctx context.Context) (lli uint32, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var doc state
	err = m.coll.FindOne(ctx, bson.M{}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = nil
		} else {
			return
		}
	}

	lli = doc.LastLoopIndex
	return
}
