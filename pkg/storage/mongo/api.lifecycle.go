package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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
