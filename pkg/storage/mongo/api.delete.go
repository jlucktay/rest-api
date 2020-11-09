package mongo

import (
	"context"

	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) Delete(id uuid.UUID) error {
	filter := bson.D{{Key: "_id", Value: mongoUUID{UUID: id}}}

	_, errDelete := s.coll.DeleteOne(context.TODO(), filter)
	if errDelete != nil {
		return errDelete
	}

	return nil
}
