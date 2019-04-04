package mongo

import (
	"github.com/jlucktay/rest-api/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

// Storage is a storage system backed by MongoDB that stores Payment structs indexed by UUID.
type Storage struct {
	client *mongo.Client
}

const (
	thisServer     = "mongodb://localhost:27017"
	thisDatabase   = "rest-api"
	thisCollection = "payments"
)

type mongoWrapper struct {
	UUID    string          `bson:"_id" json:"_id"`
	Payment storage.Payment `bson:"payment" json:"payment"`
}
