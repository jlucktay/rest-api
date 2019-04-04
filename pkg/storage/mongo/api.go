package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) Initialise() error {
	docCount, errCount := s.coll.CountDocuments(context.TODO(), bson.D{})
	if errCount != nil {
		return errCount
	}

	fmt.Printf("Collection '%s' contains %d records.\n", s.coll.Name(), docCount)

	return nil
}

func (s *Storage) Terminate() error {
	errDropColl := s.coll.Drop(context.TODO())
	if errDropColl != nil {
		return errDropColl
	}

	fmt.Printf("Collection '%s' dropped.\n", s.coll.Name())

	errDisconnect := s.coll.Database().Client().Disconnect(context.TODO())
	if errDisconnect != nil {
		return errDisconnect
	}

	fmt.Println("Connection to MongoDB closed.")

	return nil
}

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

func (s *Storage) Read(id uuid.UUID) (storage.Payment, error) {
	filter := bson.M{"_id": id.String()}

	// Create a value into which the result can be decoded.
	var found mongoWrapper
	errFind := s.coll.FindOne(context.TODO(), filter).Decode(&found)
	if errFind != nil {
		return storage.Payment{}, errFind
	}

	return found.Payment, nil
}

func (s *Storage) ReadAll(rao storage.ReadAllOptions) (map[uuid.UUID]storage.Payment, error) {
	return map[uuid.UUID]storage.Payment{}, errors.New("not yet implemented")
}

func (s *Storage) Update(id uuid.UUID, p storage.Payment) error {
	return errors.New("not yet implemented")
}

func (s *Storage) Delete(id uuid.UUID) error {
	return errors.New("not yet implemented")
}
