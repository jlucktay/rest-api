package storage_test

import (
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"

	"go.jlucktay.dev/rest-api/pkg/storage"
	"go.jlucktay.dev/rest-api/pkg/storage/inmemory"
	"go.jlucktay.dev/rest-api/pkg/storage/mongo"
	"go.jlucktay.dev/rest-api/test"
)

func TestStorage(t *testing.T) {
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
			t.Logf("Current implementation based on: %s", reflect.TypeOf(tC.ps))
			is := is.New(t)
			is.NoErr(tC.ps.Initialise())
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
			newID, errCreate := tC.ps.Create(testPayment)
			is.NoErr(errCreate)
			_, errCreateAgain := tC.ps.Create(testPayment)
			is.NoErr(errCreateAgain)

			// R
			// -> read single
			readSingle, errRead := tC.ps.Read(newID)
			is.NoErr(errRead)
			if diff := cmp.Diff(testPayment, readSingle); diff != "" {
				t.Fatalf("Mismatch (-want +got):\n%s", diff)
			}

			// -> read multiple
			var opts storage.ReadAllOptions
			readMultiple, errReadAll := tC.ps.ReadAll(opts)
			is.NoErr(errReadAll)
			is.Equal(len(readMultiple), 2) // expecting 2 records returned
			for _, actualPay := range readMultiple {
				if diff := cmp.Diff(testPayment, actualPay); diff != "" {
					t.Fatalf("Mismatch (-want +got):\n%s", diff)
				}
			}

			// U
			testPayment.Reference = "ref"
			is.NoErr(tC.ps.Update(newID, testPayment))
			updatedPay, errUpdate := tC.ps.Read(newID)
			is.NoErr(errUpdate)
			if diff := cmp.Diff(testPayment, updatedPay); diff != "" {
				t.Fatalf("Mismatch (-want +got):\n%s", diff)
			}

			// D
			is.NoErr(tC.ps.Delete(newID))
			_, errDeleted := tC.ps.Read(newID)
			is.Equal(errDeleted, &storage.NotFoundError{newID}) // expecting a NotFoundError

			// U after D
			errUpdateNonExistant := tC.ps.Update(newID, testPayment)
			is.Equal(errUpdateNonExistant, &storage.NotFoundError{newID}) // expecting a NotFoundError

			// Cleanup
			is.NoErr(tC.ps.Terminate(true))
		})
	}
}
