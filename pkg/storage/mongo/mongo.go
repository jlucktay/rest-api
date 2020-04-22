package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"

	"go.jlucktay.dev/rest-api/pkg/storage"
)

// Storage is a storage system backed by MongoDB that stores Payment structs indexed by UUID.
type Storage struct {
	coll *mongo.Collection
}

const (
	defaultServer     = "mongodb://localhost:27017"
	defaultDatabase   = "rest-api"
	defaultCollection = "payments"
)

type mongoWrapper struct {
	UUID    mongoUUID       `bson:"_id" json:"_id"`
	Payment storage.Payment `bson:"payment" json:"payment"`
}

type OptionEnum uint

const (
	Server OptionEnum = iota
	Database
	Collection
)

type Option struct {
	Key   OptionEnum
	Value string
}
