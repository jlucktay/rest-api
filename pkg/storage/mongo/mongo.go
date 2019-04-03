package mongo

import (
	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
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
	UUID    *mongoUUID      `bson:"uuid" json:"uuid"`
	Payment storage.Payment `bson:"payment" json:"payment"`
}

type mongoUUID struct {
	uuid.UUID
}

func (mu *mongoUUID) MarshalBSON() ([]byte, error) {
	return []byte(mu.UUID.String()), nil
}

func (mu *mongoUUID) UnmarshalBSON(b []byte) error {
	mu.UUID = uuid.FromStringOrNil(string(b))
	return nil
}
