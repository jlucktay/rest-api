package mongo

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jlucktay/rest-api/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) Update(id uuid.UUID, p storage.Payment) error {
	filter := bson.D{{Key: "_id", Value: mongoUUID(id)}}
	mongoUpdate := bson.M{"$set": wrap(p, mongoUUID(id))}
	_, errUpdate := s.coll.UpdateOne(context.TODO(), filter, mongoUpdate)
	if errUpdate != nil {
		return errUpdate
		// return &storage.NotFoundError{ID: id} // todo?
	}
	return nil
}