package mongo

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) Delete(id uuid.UUID) error {
	filter := bson.D{{Key: "_id", Value: id.String()}}
	_, errDelete := s.coll.DeleteOne(context.TODO(), filter)
	if errDelete != nil {
		return errDelete
	}
	return nil
}
