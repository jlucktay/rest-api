package mongo

import (
	"context"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) Update(id uuid.UUID, p storage.Payment) error {
	filter := bson.D{{Key: "_id", Value: id.String()}}
	mongoUpdate := bson.M{"$set": wrap(p, id)}
	_, errUpdate := s.coll.UpdateOne(context.TODO(), filter, mongoUpdate)
	if errUpdate != nil {
		return errUpdate
		// return &storage.NotFoundError{ID: id} // todo?
	}
	return nil
}
