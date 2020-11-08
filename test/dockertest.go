package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupDocker sets up a MongoDB container with dockertest and returns a purge function (which should be deferred)
// and a connection string.
func SetupDocker(t *testing.T) (purge func(), cs string) {
	is := is.New(t)
	// Set up a disposable MongoDB container with Docker
	t.Log("Docker/MongoDB starting...")

	dockerPool, errPool := dockertest.NewPool("")
	is.NoErr(errPool)

	mongoResource, errRun := dockerPool.Run("mongo", "4", nil)
	is.NoErr(errRun)

	mongoContainerName := mongoResource.Container.Name

	purge = func() {
		t.Logf("Purging Docker/MongoDB container '%s'...", mongoContainerName)
		is.NoErr(dockerPool.Purge(mongoResource))
		t.Logf("Purged Docker/MongoDB container '%s'!", mongoContainerName)
	}

	cs = fmt.Sprintf(
		"mongodb://localhost:%s/jra_test",
		mongoResource.GetPort("27017/tcp"))

	// Exponential backoff-retry, while MongoDB gets ready to accept connections
	if err := dockerPool.Retry(func() error {
		mgoOpts := (&options.ClientOptions{}).ApplyURI(cs)

		mgoClient, errConnect := mongo.Connect(context.TODO(), mgoOpts)
		if errConnect != nil {
			return errConnect
		}

		return mgoClient.Ping(context.TODO(), nil)
	}); err != nil {
		is.NoErr(err) // could not ping Docker/MongoDB
	}

	t.Logf("Started Docker/MongoDB container '%s'!", mongoContainerName)

	return purge, cs
}
