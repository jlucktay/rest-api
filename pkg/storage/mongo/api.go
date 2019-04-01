package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) Initialise() error {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	var errConnect error
	s.client, errConnect = mongo.Connect(context.TODO(), clientOptions)

	if errConnect != nil {
		return errConnect
	}

	// Check the connection
	errPing := s.client.Ping(context.TODO(), nil)

	if errPing != nil {
		return errPing
	}

	fmt.Println("Connected to MongoDB!")

	collection := s.client.Database(thisDatabase).Collection(thisCollection)
	docCount, errCount := collection.CountDocuments(context.TODO(), bson.D{})
	if errCount != nil {
		return errCount
	}

	fmt.Printf("Collection '%s' contains %d records.\n", collection.Name(), docCount)

	return nil
}

func (s *Storage) Terminate() error {
	err := s.client.Disconnect(context.TODO())

	if err != nil {
		return err
	}

	fmt.Println("Connection to MongoDB closed.")

	return nil
}

func (s *Storage) Create(newPayment storage.Payment) (uuid.UUID, error) {
	newID, errNew := uuid.NewV4()
	if errNew != nil {
		return uuid.Nil, errNew
	}

	mongoInsert := bson.M{
		"_id":     newID.String(),
		"payment": newPayment,
	}

	c := s.client.Database(thisDatabase).Collection(thisCollection)

	insertResult, errInsert := c.InsertOne(context.TODO(), mongoInsert)
	if errInsert != nil {
		return uuid.Nil, errInsert
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return newID, nil
}

func (s *Storage) CreateSpecificID(newID uuid.UUID, newPayment storage.Payment) error {
	mongoInsert := bson.M{
		"_id":     newID.String(),
		"payment": newPayment,
	}

	c := s.client.Database(thisDatabase).Collection(thisCollection)

	insertResult, errInsert := c.InsertOne(context.TODO(), mongoInsert)
	if errInsert != nil {
		return errInsert
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return nil
}

func (s *Storage) Read(id uuid.UUID) (storage.Payment, error) {
	return storage.Payment{}, errors.New("not yet implemented")
}

func (s *Storage) ReadAll(rao storage.ReadAllOptions) (map[uuid.UUID]storage.Payment, error) {
	return make(map[uuid.UUID]storage.Payment), errors.New("not yet implemented")
}

func (s *Storage) Update(id uuid.UUID, p storage.Payment) error {
	return errors.New("not yet implemented")
}

func (s *Storage) Delete(id uuid.UUID) error {
	return errors.New("not yet implemented")
}
