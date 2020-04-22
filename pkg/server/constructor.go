package server

import (
	"os"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"go.jlucktay.dev/rest-api/pkg/storage/inmemory"
	"go.jlucktay.dev/rest-api/pkg/storage/mongo"
)

// New creates a new Server utilising the given StorageType to handle Payment storage, and sets up the HTTP router.
// It takes an optional string argument which will set the hostname of the MongoDB server to connect to.
func New(st StorageType, logDebug bool, host ...string) *Server {
	s := &Server{Router: chi.NewRouter()}
	s.setupRoutes()

	switch st {
	case InMemory:
		s.Storage = &inmemory.Storage{}
	case Mongo:
		mongoServer := defaultServer
		if len(host) > 0 {
			mongoServer = "mongodb://" + host[0] + ":27017"
		}

		s.Storage = mongo.New(
			mongo.Option{
				Key:   mongo.Server,
				Value: mongoServer,
			},
			mongo.Option{
				Key:   mongo.Database,
				Value: defaultDatabase,
			},
			mongo.Option{
				Key:   mongo.Collection,
				Value: defaultCollection,
			},
		)
	}

	log.SetLevel(log.ErrorLevel)
	log.SetOutput(os.Stderr)
	log.SetFormatter(new(log.JSONFormatter))

	if logDebug {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		log.Debug("Debug logging is enabled.")
	}

	if errStorageInit := s.Storage.Initialise(); errStorageInit != nil {
		log.Fatal(errStorageInit)
	}

	return s
}
