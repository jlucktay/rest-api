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

func (s *Storage) Read(id uuid.UUID) (storage.Payment, error) {
	filter := bson.M{"_id": mongoUUID(id)}

	// Create a value into which the result can be decoded.
	var found mongoWrapper
	errFind := s.coll.FindOne(context.TODO(), filter).Decode(&found)
	if errFind != nil {
		if errFind.Error() == "mongo: no documents in result" {
			return storage.Payment{}, &storage.NotFoundError{ID: id}
		}
		return storage.Payment{}, errFind
	}

	return found.Payment, nil
}

func (s *Storage) ReadAll(rao storage.ReadAllOptions) (map[uuid.UUID]storage.Payment, error) {
	// Set limit from options or default constant.
	if rao.Limit == 0 {
		rao.Limit = storage.DefaultLimit
	}

	// Get all keys and sort in order.
	filter := bson.D{} // #nofilter

	opts := &options.FindOptions{}  // Sort UUIDs ascending.
	opts.SetLimit(int64(rao.Limit)) // No fear of losng data when casting like this, as they are both originally uint.
	opts.SetSkip(int64(rao.Offset))
	opts.SetSort(bson.D{{Key: "_id", Value: 1}})

	cur, errFind := s.coll.Find(context.TODO(), filter, opts)
	if errFind != nil {
		return nil, fmt.Errorf("couldn't find records with given parameters: %v", errFind)
	}

	defer cur.Close(context.TODO())
	found := make(map[uuid.UUID]storage.Payment)

	for cur.Next(context.TODO()) {
		mwDec := mongoWrapper{}
		if errDecode := cur.Decode(&mwDec); errDecode != nil {
			return nil, fmt.Errorf("couldn't make element ready for display: %v", errDecode)
		}

		found[uuid.UUID(mwDec.UUID)] = mwDec.Payment // Unwrap
	}
	if cur.Err() != nil {
		return nil, errors.New("cursor error")
	}

	return found, nil
}
