package state

import (
	"github.com/leaq-ru/org-producer/mongo"
	md "go.mongodb.org/mongo-driver/mongo"
)

func NewModel(db *md.Database) Model {
	return Model{
		coll: db.Collection(mongo.CollState),
	}
}
