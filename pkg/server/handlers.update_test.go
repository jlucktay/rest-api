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
)

func TestUpdatePayment(t *testing.T) {
	s := server.New(server.InMemory, false)
	var w *httptest.ResponseRecorder
	i := is.New(t)

	startAmount := 123.45
	updatedAmount := 246.80

	// Construct a HTTP request which creates a payment.
	p := storage.Payment{Amount: startAmount}
	j, errMarshalCreate := json.Marshal(p)
	i.NoErr(errMarshalCreate)
	reqBodyCreate := bytes.NewBuffer(j)
	reqCreate, errCreate := http.NewRequest(http.MethodPost, "/v1/payments", reqBodyCreate)
	i.NoErr(errCreate)
	reqCreate.Header.Set("Content-Type", "application/json")

	// Send it, and record the HTTP back and forth.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqCreate)
	resp := w.Result()
	defer resp.Body.Close()
	i.Equal(http.StatusCreated, resp.StatusCode)

	// Get the Location header which points at the new payment.
	loc := resp.Header.Get("Location")
	r := regexp.MustCompile("^/v1/payments/([0-9a-f-]{36})$")
	i.True(r.MatchString(loc))
	newID := r.FindStringSubmatch(loc)[1]

	// Construct another HTTP request to update the new payment.
	p.Amount = updatedAmount
	k, errMarshalUpdate := json.Marshal(p)
	i.NoErr(errMarshalUpdate)
	reqBodyUpdate := bytes.NewBuffer(k)
	reqUpdate, errUpdate := http.NewRequest(http.MethodPut, "/v1/payments/"+newID, reqBodyUpdate)
	i.NoErr(errUpdate)

	// Update the payment using the ID returned via 'Location' header.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqUpdate)
	resp2 := w.Result()
	defer resp2.Body.Close()
	i.Equal(http.StatusNoContent, resp2.StatusCode)

	// Construct another HTTP request to read the payment.
	reqRead, errRead := http.NewRequest(http.MethodGet, "/v1/payments/"+newID, nil)
	i.NoErr(errRead)

	// Send the read request and assert on the length of the response.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqRead)
	resp3 := w.Result()
	defer resp3.Body.Close()
	i.Equal(http.StatusOK, resp3.StatusCode)
	respBodyBytes, errReadAll := ioutil.ReadAll(resp3.Body)
	i.NoErr(errReadAll)
	i.True(len(string(respBodyBytes)) > 0)

	// Unmarshal into a slice of Payment structs.
	returnedPayment := server.NewWrapper("")
	errUnmarshal := json.Unmarshal(respBodyBytes, &returnedPayment)
	i.NoErr(errUnmarshal)

	// Assert that the Amount of the Payment returned has been updated appropriately.
	i.Equal(updatedAmount, returnedPayment.Data[0].Attributes.Amount)
}
