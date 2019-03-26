package main

import (
	"reflect"
	"testing"

	"github.com/matryer/is"
	"github.com/shopspring/decimal"
)

func TestStorage(t *testing.T) {
	testCases := []struct {
		desc string
		ps   PaymentStorage
	}{
		{
			desc: "In-memory storage (map); won't persist across app restarts",
			ps:   &inMemoryStorage{},
		},
	}
	for _, tC := range testCases {
		tC := tC // pin!
		t.Run(tC.desc, func(t *testing.T) {
			t.Logf("Current implementation based on: %s", reflect.TypeOf(tC.ps))
			i := is.New(t)
			i.NoErr(tC.ps.Init())
			testPayment := Payment{Amount: decimal.NewFromFloat(123.45)}

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
			var opts ReadAllOptions
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
