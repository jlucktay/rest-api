package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) Initialise() error {
	filter := bson.D{} // #nofilter
	docCount, errCount := s.coll.CountDocuments(context.TODO(), filter)
	if errCount != nil {
		return errCount
	}

	fmt.Printf("Collection '%s' contains %d records.\n", s.coll.Name(), docCount)

	return nil
}

// Terminate will keep or destroy data, depending on whether true or false is passed in.
// Will default to keeping data if no bool is given.
func (s *Storage) Terminate(drop ...bool) error {
	if len(drop) > 0 && drop[0] {
		errDropColl := s.coll.Drop(context.TODO())
		if errDropColl != nil {
			return errDropColl
		}

		fmt.Printf("Collection '%s' dropped.\n", s.coll.Name())
	}

	errDisconnect := s.coll.Database().Client().Disconnect(context.TODO())
	if errDisconnect != nil {
		return errDisconnect
	}

	fmt.Println("Disconnected from MongoDB.")

	return nil
}
