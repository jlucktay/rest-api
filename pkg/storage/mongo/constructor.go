package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(mo ...Option) *Storage {
	server := defaultServer
	database := defaultDatabase
	collection := defaultCollection

	for _, opt := range mo {
		switch opt.Key {
		case Server:
			server = opt.Value
		case Database:
			database = opt.Value
		case Collection:
			collection = opt.Value
		}
	}

	s := &Storage{
		coll: connect(server).Database(database).Collection(collection),
	}

	return s
}

func connect(server string) *mongo.Client {
	// Set client options.
	clientOptions := options.Client().ApplyURI(server)

	// Connect to MongoDB.
	client, errConnect := mongo.Connect(context.TODO(), clientOptions)

	if errConnect != nil {
		log.Fatal(errConnect)
	}

	// Check the connection.
	errPing := client.Ping(context.TODO(), nil)

	if errPing != nil {
		log.Fatal(errPing)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}
