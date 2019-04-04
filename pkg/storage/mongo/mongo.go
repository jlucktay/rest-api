package mongo

import (
	"github.com/jlucktay/rest-api/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
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
	UUID    string          `bson:"_id" json:"_id"`
	Payment storage.Payment `bson:"payment" json:"payment"`
}

type MongoOptionEnum uint

const (
	Server MongoOptionEnum = iota
	Database
	Collection
)

type MongoOption struct {
	Key   MongoOptionEnum
	Value string
}
