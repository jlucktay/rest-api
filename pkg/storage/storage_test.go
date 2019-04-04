package storage_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/jlucktay/rest-api/pkg/storage/inmemory"
	"github.com/jlucktay/rest-api/pkg/storage/mongo"
	"github.com/matryer/is"
)

func TestStorage(t *testing.T) {
	testCases := []struct {
		desc string
		ps   storage.PaymentStorage
	}{
		{
			desc: "In-memory storage (map); won't persist across app restarts",
			ps:   &inmemory.Storage{},
		},
		{
			desc: "Database storage (MongoDB); will persist across app restarts",
			ps:   &mongo.Storage{},
		},
	}
	for _, tC := range testCases {
		tC := tC // pin!
		t.Run(tC.desc, func(t *testing.T) {
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
				i.True(reflect.DeepEqual(testPayment, actualPay))
			}

			// U
			testPayment.Reference = "ref"
			i.NoErr(tC.ps.Update(newID, testPayment))
			updatedPay, errUpdate := tC.ps.Read(newID)
			i.NoErr(errUpdate)
			i.True(reflect.DeepEqual(testPayment, updatedPay))

			// D
			i.NoErr(tC.ps.Delete(newID))
			_, errDeleted := tC.ps.Read(newID)
			i.Equal(errDeleted, &storage.NotFoundError{newID})

			// Cleanup
			i.NoErr(tC.ps.Terminate())
		})
	}
}
