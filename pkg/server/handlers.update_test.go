package server_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/matryer/is"
	"github.com/shopspring/decimal"
)

func TestUpdatePayment(t *testing.T) {
	s := server.New(server.InMemory)
	var w *httptest.ResponseRecorder
	i := is.New(t)

	startAmount := decimal.NewFromFloat(123.45)
	updatedAmount := decimal.NewFromFloat(246.80)

	// Construct a HTTP request which creates a payment.
	p := storage.Payment{Amount: startAmount}
	j, errMarshalCreate := json.Marshal(p)
	i.NoErr(errMarshalCreate)
	reqBodyCreate := bytes.NewBuffer(j)
	reqCreate, errCreate := http.NewRequest(http.MethodPost, "/payments", reqBodyCreate)
	i.NoErr(errCreate)
	reqCreate.Header.Set("Content-Type", "application/json")

	// Send it, and record the HTTP back and forth.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqCreate)
	i.Equal(http.StatusCreated, w.Result().StatusCode)

	// Get the Location header which points at the new payment.
	loc := w.Result().Header.Get("Location")
	r := regexp.MustCompile("^/payments/([0-9a-f-]{36})$")
	i.True(r.MatchString(loc))
	newID := r.FindStringSubmatch(loc)[1]

	// Construct another HTTP request to update the new payment.
	p.Amount = updatedAmount
	k, errMarshalUpdate := json.Marshal(p)
	i.NoErr(errMarshalUpdate)
	reqBodyUpdate := bytes.NewBuffer(k)
	reqUpdate, errUpdate := http.NewRequest(http.MethodPut, "/payments/"+newID, reqBodyUpdate)
	i.NoErr(errUpdate)

	// Update the payment using the ID returned via 'Location' header.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqUpdate)
	i.Equal(http.StatusNoContent, w.Result().StatusCode)

	// Construct another HTTP request to read the payment.
	reqRead, errRead := http.NewRequest(http.MethodGet, "/payments/"+newID, nil)
	i.NoErr(errRead)

	// Send the read request and assert on the length of the response.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqRead)
	i.Equal(http.StatusOK, w.Result().StatusCode)
	respBodyBytes, errReadAll := ioutil.ReadAll(w.Result().Body)
	i.NoErr(errReadAll)
	i.True(len(string(respBodyBytes)) > 0)

	// Unmarshal into a slice of Payment structs.
	var returnedPayment server.ReadWrapper
	errUnmarshal := json.Unmarshal(respBodyBytes, &returnedPayment)
	i.NoErr(errUnmarshal)

	// Assert that the Amount of the Payment returned has been updated appropriately.
	i.Equal(updatedAmount, returnedPayment.Data[0].Attributes.Amount)
}
