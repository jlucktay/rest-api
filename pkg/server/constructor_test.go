package server_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/matryer/is"
	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.jlucktay.dev/rest-api/pkg/server"
	"go.jlucktay.dev/rest-api/pkg/storage"
	"go.jlucktay.dev/rest-api/pkg/storage/inmemory"
	jramongo "go.jlucktay.dev/rest-api/pkg/storage/mongo"
)

func TestNew(t *testing.T) {
	is := is.New(t)

	// Set up a disposable MongoDB container with Docker
	t.Log("Docker/MongoDB starting...")

	dockerPool, errPool := dockertest.NewPool("")
	is.NoErr(errPool)

	mongoResource, errRun := dockerPool.Run("mongo", "4", nil)
	is.NoErr(errRun)

	mongoContainerName := mongoResource.Container.Name

	defer func() {
		t.Logf("Purging Docker/MongoDB container '%s'...", mongoContainerName)
		is.NoErr(dockerPool.Purge(mongoResource))
		t.Logf("Purged Docker/MongoDB container '%s'!", mongoContainerName)
	}()

	directConnString := fmt.Sprintf(
		"mongodb://localhost:%s/jra_test",
		mongoResource.GetPort("27017/tcp"))

	// Exponential backoff-retry, while MongoDB gets ready to accept connections
	if err := dockerPool.Retry(func() error {
		mgoOpts := (&options.ClientOptions{}).ApplyURI(directConnString)

		mgoClient, errConnect := mongo.Connect(context.TODO(), mgoOpts)
		if errConnect != nil {
			return errConnect
		}

		return mgoClient.Ping(context.TODO(), nil)
	}); err != nil {
		is.NoErr(err) // could not ping Docker/MongoDB
	}

	t.Logf("Started Docker/MongoDB container '%s'!", mongoContainerName)

	testCases := map[string]struct {
		st         server.StorageType
		expected   storage.PaymentStorage
		connString string
	}{
		"In-memory storage (map); won't persist across app restarts": {
			st:       server.InMemory,
			expected: &inmemory.Storage{},
		},
		"Database storage (MongoDB); will persist across app restarts": {
			st:         server.Mongo,
			expected:   &jramongo.Storage{},
			connString: directConnString,
		},
	}

	for name, tC := range testCases {
		name, tC := name, tC // pin!

		t.Run(name, func(_ *testing.T) {
			var s *server.Server
			if tC.connString != "" {
				s = server.New(tC.st, true, tC.connString)
			} else {
				s = server.New(tC.st, true)
			}
			is.Equal(reflect.TypeOf(s.Storage), reflect.TypeOf(tC.expected)) // wrong storage type
		})
	}
}
