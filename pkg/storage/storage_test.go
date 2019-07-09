package storage_test

import (
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/jlucktay/rest-api/pkg/storage/inmemory"
	"github.com/jlucktay/rest-api/pkg/storage/mongo"
	"github.com/matryer/is"
)

func TestStorage(t *testing.T) {
	t.Parallel() // parallelise with other tests

	randTestID := uuid.Must(uuid.NewV4())

	testCases := map[string]struct {
		ps storage.PaymentStorage
	}{
		"In-memory storage (map); won't persist across app restarts": {
			ps: &inmemory.Storage{},
		},
		"Database storage (MongoDB); will persist across app restarts": {
			ps: mongo.New(
				mongo.Option{
					Key:   mongo.Database,
					Value: "test",
				},
				mongo.Option{
					Key:   mongo.Collection,
					Value: "test-" + randTestID.String(),
				},
			),
		},
	}
	for name, tC := range testCases {
		tC := tC // pin!
		t.Run(name, func(t *testing.T) {
			t.Parallel() // parallelise with other sub-tests

			t.Logf("Current implementation based on: %s", reflect.TypeOf(tC.ps))
			i := is.New(t)
			i.NoErr(tC.ps.Initialise())
			testPayment := storage.Payment{
				Amount: 123.45,
				ChargesInformation: storage.ChargesInformation{
					SenderCharges: []storage.SenderCharges{
						{Amount: 1.01},
						{Amount: 2.02},
						{Amount: 3.03},
					},
				},
				PaymentID: "test",
			}

			// C
			newID, errCreate := tC.ps.Create(testPayment)
			i.NoErr(errCreate)
			_, errCreateAgain := tC.ps.Create(testPayment)
			i.NoErr(errCreateAgain)

			// R
			// -> read single
			readSingle, errRead := tC.ps.Read(newID)
			i.NoErr(errRead)
			if diff := cmp.Diff(testPayment, readSingle); diff != "" {
				t.Fatalf("Mismatch (-want +got):\n%s", diff)
			}

			// -> read multiple
			var opts storage.ReadAllOptions
			readMultiple, errReadAll := tC.ps.ReadAll(opts)
			i.NoErr(errReadAll)
			i.Equal(len(readMultiple), 2)
			for _, actualPay := range readMultiple {
				if diff := cmp.Diff(testPayment, actualPay); diff != "" {
					t.Fatalf("Mismatch (-want +got):\n%s", diff)
				}
			}

			// U
			testPayment.Reference = "ref"
			i.NoErr(tC.ps.Update(newID, testPayment))
			updatedPay, errUpdate := tC.ps.Read(newID)
			i.NoErr(errUpdate)
			if diff := cmp.Diff(testPayment, updatedPay); diff != "" {
				t.Fatalf("Mismatch (-want +got):\n%s", diff)
			}

			// D
			i.NoErr(tC.ps.Delete(newID))
			_, errDeleted := tC.ps.Read(newID)
			i.Equal(errDeleted, &storage.NotFoundError{newID})

			// U after D
			errUpdateNonExistant := tC.ps.Update(newID, testPayment)
			i.Equal(errUpdateNonExistant, &storage.NotFoundError{newID})

			// Cleanup
			i.NoErr(tC.ps.Terminate(true))
		})
	}
}
