package server

import (
	"log"

	"github.com/go-chi/chi"
	"github.com/jlucktay/rest-api/pkg/storage/inmemory"
	"github.com/jlucktay/rest-api/pkg/storage/mongo"
)

// New creates a new Server utilising the given StorageType to handle Payment storage, and sets up the HTTP router.
func New(st StorageType) *Server {
	s := &Server{Router: chi.NewRouter()}
	s.setupRoutes()

	switch st {
	case InMemory:
		s.Storage = &inmemory.Storage{}
	case Mongo:
		s.Storage = mongo.New(
			mongo.Option{
				Key:   mongo.Database,
				Value: "rest-api",
			},
			mongo.Option{
				Key:   mongo.Collection,
				Value: "payments",
			},
		)
	}

	if errStorageInit := s.Storage.Initialise(); errStorageInit != nil {
		log.Fatal(errStorageInit)
	}

	return s
}
