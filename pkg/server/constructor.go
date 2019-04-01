package server

import (
	"log"

	"github.com/jlucktay/rest-api/pkg/storage/inmemory"
	"github.com/julienschmidt/httprouter"
)

// New creates a new Server utilising the given StorageType to handle Payment storage, and sets up the HTTP router.
func New(st StorageType) (s *Server) {
	s = &Server{
		Router: httprouter.New(),
	}
	s.setupRoutes()

	switch st {
	case InMemory:
		s.Storage = &inmemory.Storage{}
	case Mongo:
		panic("Mongo storage is not yet implemented")
	}

	if errStorageInit := s.Storage.Initialise(); errStorageInit != nil {
		log.Fatal(errStorageInit)
	}

	return
}
