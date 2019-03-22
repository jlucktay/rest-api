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

			// R
			readPay, errRead := tC.ps.Read(newID)
			i.NoErr(errRead)
			i.True(reflect.DeepEqual(testPayment, readPay))
			_, _ = tC.ps.Create(testPayment)
			var opts ReadAllOptions
			readAllPay, errReadAll := tC.ps.ReadAll(opts)
			i.NoErr(errReadAll)
			i.Equal(len(readAllPay), 2)
			i.True(reflect.DeepEqual(testPayment, readAllPay[0]))
			i.True(reflect.DeepEqual(testPayment, readAllPay[1]))

			// U
			testPayment.Reference = "ref"
			i.NoErr(tC.ps.Update(newID, testPayment))
			updatedPay, _ := tC.ps.Read(newID)
			i.True(reflect.DeepEqual(testPayment, updatedPay))

			// D
			i.NoErr(tC.ps.Delete(newID))
			_, errDeleted := tC.ps.Read(newID)
			i.Equal(errDeleted, &NotFoundError{newID})
		})
	}
}
