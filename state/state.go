package state

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type state struct {
	ID            primitive.ObjectID `bson:"_id"`
	LastLoopIndex uint32             `bson:"l"`
	UpdatedAt     time.Time          `bson:"u"`
}
