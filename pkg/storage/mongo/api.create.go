package mongo

import (
	"context"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

func (s *Storage) Create(newPayment storage.Payment) (uuid.UUID, error) {
	mongoInsert := wrap(newPayment)

	_, errInsert := s.coll.InsertOne(context.TODO(), mongoInsert)
	if errInsert != nil {
		return uuid.Nil, errInsert
	}

	return uuid.FromStringOrNil(mongoInsert.UUID), nil
}

func (s *Storage) CreateSpecificID(newID uuid.UUID, newPayment storage.Payment) error {
	mongoInsert := wrap(newPayment, newID)

	_, errInsert := s.coll.InsertOne(context.TODO(), mongoInsert)
	if errInsert != nil {
		return errInsert
	}

	return nil
}
