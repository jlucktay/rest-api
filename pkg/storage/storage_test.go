package storage_test

import (
	"reflect"
	"testing"

	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/jlucktay/rest-api/pkg/storage/inmemory"
	"github.com/matryer/is"
	"github.com/shopspring/decimal"
)

func TestStorage(t *testing.T) {
	testCases := []struct {
		desc string
		ps   storage.PaymentStorage
	}{
		{
			desc: "In-memory storage (map); won't persist across app restarts",
			ps:   &inmemory.InMemoryStorage{},
		},
	}
	for _, tC := range testCases {
		tC := tC // pin!
		t.Run(tC.desc, func(t *testing.T) {
			t.Logf("Current implementation based on: %s", reflect.TypeOf(tC.ps))
			i := is.New(t)
			i.NoErr(tC.ps.Init())
			testPayment := storage.Payment{Amount: decimal.NewFromFloat(123.45)}

			// C
			newID, errCreate := tC.ps.Create(testPayment)
			i.NoErr(errCreate)
			_, errCreateAgain := tC.ps.Create(testPayment)
			i.NoErr(errCreateAgain)

			// R
			// -> read single
			readSingle, errRead := tC.ps.Read(newID)
			i.NoErr(errRead)
			i.True(reflect.DeepEqual(testPayment, readSingle))

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
			i.Equal(errDeleted, &NotFoundError{newID})
		})
	}
}
