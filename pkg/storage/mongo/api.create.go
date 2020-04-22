package mongo

import (
	"context"

	"github.com/gofrs/uuid"

	"go.jlucktay.dev/rest-api/pkg/storage"
)

func (s *Storage) Create(newPayment storage.Payment) (uuid.UUID, error) {
	mongoInsert := wrap(newPayment)

	_, errInsert := s.coll.InsertOne(context.TODO(), mongoInsert)
	if errInsert != nil {
		return uuid.Nil, errInsert
	}

	return mongoInsert.UUID.UUID, nil
}

func (s *Storage) CreateSpecificID(newID uuid.UUID, newPayment storage.Payment) error {
	mongoInsert := wrap(newPayment, mongoUUID{newID})

	_, errInsert := s.coll.InsertOne(context.TODO(), mongoInsert)
	if errInsert != nil {
		return errInsert
	}

	return nil
}
