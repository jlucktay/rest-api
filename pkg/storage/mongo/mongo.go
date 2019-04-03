package mongo

import (
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
