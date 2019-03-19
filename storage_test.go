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
		t.Run(tC.desc, func(t *testing.T) {
			t.Logf("Current implementation based on: %s", reflect.TypeOf(tC.ps))
			i := is.New(t)
			i.NoErr(tC.ps.Init())
			testPayment := Payment{Amount: decimal.NewFromFloat(123.45)}

			// C
			newId, errCreate := tC.ps.Create(testPayment)
			i.NoErr(errCreate)

			// R
			readPay, errRead := tC.ps.Read(newId)
			i.NoErr(errRead)
			i.True(reflect.DeepEqual(testPayment, readPay))
			_, _ = tC.ps.Create(testPayment)
			readAllPay, errReadAll := tC.ps.ReadAll()
			i.NoErr(errReadAll)
			i.Equal(len(readAllPay), 2)
			i.True(reflect.DeepEqual(testPayment, readAllPay[0]))
			i.True(reflect.DeepEqual(testPayment, readAllPay[1]))

			// U
			testPayment.Reference = "ref"
			i.NoErr(tC.ps.Update(newId, testPayment))
			updatedPay, _ := tC.ps.Read(newId)
			i.True(reflect.DeepEqual(testPayment, updatedPay))

			// D
			i.NoErr(tC.ps.Delete(newId))
			_, errDeleted := tC.ps.Read(newId)
			i.Equal(errDeleted, &NotFoundError{newId})
		})
	}
}
