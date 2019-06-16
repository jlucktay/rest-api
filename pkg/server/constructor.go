package server

import (
	"github.com/go-chi/chi"
	"github.com/jlucktay/rest-api/pkg/storage/inmemory"
	"github.com/jlucktay/rest-api/pkg/storage/mongo"
	log "github.com/sirupsen/logrus"
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

	s.Log = &logrus.Logger{
		Formatter: new(logrus.JSONFormatter),
		Level:     logrus.ErrorLevel,
		Out:       os.Stderr,
	}

	if logDebug {
		s.Log.SetLevel(logrus.DebugLevel)
		s.Log.SetReportCaller(true)
		s.Log.Debug("Debug logging is enabled.")
	}

	if errStorageInit := s.Storage.Initialise(); errStorageInit != nil {
		log.Fatal(errStorageInit)
	}

	return s
}
