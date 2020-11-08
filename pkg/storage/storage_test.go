package storage_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"
	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.jlucktay.dev/rest-api/pkg/storage"
	"go.jlucktay.dev/rest-api/pkg/storage/inmemory"
	jramongo "go.jlucktay.dev/rest-api/pkg/storage/mongo"
	"go.jlucktay.dev/rest-api/test"
)

func TestStorage(t *testing.T) { //nolint:funlen
	is := is.New(t)
	randTestID := uuid.Must(uuid.NewV4())

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

	testCases := map[string]storage.PaymentStorage{
		"In-memory storage (map); won't persist across app restarts": &inmemory.Storage{},

		"Database storage (MongoDB); will persist across app restarts": jramongo.New(
			jramongo.Option{Key: jramongo.Server, Value: directConnString},
			jramongo.Option{Key: jramongo.Database, Value: "test"},
			jramongo.Option{Key: jramongo.Collection, Value: "test-" + randTestID.String()},
		),
	}

	for name, tC := range testCases {
		// pin! ref: https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		name, tC := name, tC

		t.Run(name, func(t *testing.T) {
			// t.Parallel() // Don't use .Parallel() without pinning.

			t.Logf("Current implementation based on: %s", reflect.TypeOf(tC))
			is.NoErr(tC.Initialise())
			testPayment := storage.Payment{
				Amount: test.Amount,
				ChargesInformation: storage.ChargesInformation{
					SenderCharges: []storage.SenderCharges{
						{Amount: test.ChargeOne},
						{Amount: test.ChargeTwo},
						{Amount: test.ChargeThree},
					},
				},
				PaymentID: "test",
			}

			// C
			newID, errCreate := tC.Create(testPayment)
			is.NoErr(errCreate)
			_, errCreateAgain := tC.Create(testPayment)
			is.NoErr(errCreateAgain)

			// R
			// -> read single
			readSingle, errRead := tC.Read(newID)
			is.NoErr(errRead)
			if diff := cmp.Diff(testPayment, readSingle); diff != "" {
				t.Fatalf("Mismatch (-want +got):\n%s", diff)
			}

			// -> read multiple
			var opts storage.ReadAllOptions
			readMultiple, errReadAll := tC.ReadAll(opts)
			is.NoErr(errReadAll)
			is.Equal(len(readMultiple), 2) // expecting 2 records returned
			for _, actualPay := range readMultiple {
				if diff := cmp.Diff(testPayment, actualPay); diff != "" {
					t.Fatalf("Mismatch (-want +got):\n%s", diff)
				}
			}

			// U
			testPayment.Reference = "ref"
			is.NoErr(tC.Update(newID, testPayment))
			updatedPay, errUpdate := tC.Read(newID)
			is.NoErr(errUpdate)
			if diff := cmp.Diff(testPayment, updatedPay); diff != "" {
				t.Fatalf("Mismatch (-want +got):\n%s", diff)
			}

			// D
			is.NoErr(tC.Delete(newID))
			_, errDeleted := tC.Read(newID)
			is.Equal(errDeleted, &storage.NotFoundError{newID}) // expecting a NotFoundError

			// U after D
			errUpdateNonExistant := tC.Update(newID, testPayment)
			is.Equal(errUpdateNonExistant, &storage.NotFoundError{newID}) // expecting a NotFoundError

			// Cleanup
			is.NoErr(tC.Terminate(true))
		})
	}
}
