package server_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/rest-api/pkg/server"
	"go.jlucktay.dev/rest-api/pkg/storage"
)

func checkResponseStatus(t *testing.T, s *server.Server, r *http.Request, expCode int) *http.Response {
	t.Helper()

	w := httptest.NewRecorder()
	is := is.New(t)

	s.Router.ServeHTTP(w, r)
	resp := w.Result()
	is.Equal(expCode, resp.StatusCode) // expecting HTTP 201

	return resp
}

func TestUpdatePayment(t *testing.T) {
	s := server.New(server.InMemory, false)
	is := is.New(t)

	startAmount := 123.45
	updatedAmount := 246.80

	// Construct a HTTP request which creates a payment.
	p := storage.Payment{Amount: startAmount}
	j, errMarshalCreate := json.Marshal(p)
	is.NoErr(errMarshalCreate)

	reqBodyCreate := bytes.NewBuffer(j)
	reqCreate, errCreate := http.NewRequest(http.MethodPost, "/v1/payments", reqBodyCreate)

	is.NoErr(errCreate)
	reqCreate.Header.Set("Content-Type", "application/json")

	// Send it, and record the HTTP back and forth.
	resp1 := checkResponseStatus(t, s, reqCreate, http.StatusCreated) // expecting HTTP 201
	defer resp1.Body.Close()

	// Get the Location header which points at the new payment.
	loc := resp1.Header.Get("Location")
	r := regexp.MustCompile("^/v1/payments/([0-9a-f-]{36})$")
	is.True(r.MatchString(loc)) // regex couldn't match Location header
	newID := r.FindStringSubmatch(loc)[1]

	// Construct another HTTP request to update the new payment.
	p.Amount = updatedAmount
	k, errMarshalUpdate := json.Marshal(p)
	is.NoErr(errMarshalUpdate)

	reqBodyUpdate := bytes.NewBuffer(k)
	reqUpdate, errUpdate := http.NewRequest(http.MethodPut, "/v1/payments/"+newID, reqBodyUpdate)

	is.NoErr(errUpdate)

	// Update the payment using the ID returned via 'Location' header.
	resp2 := checkResponseStatus(t, s, reqUpdate, http.StatusNoContent) // expecting HTTP 204
	defer resp2.Body.Close()

	// Construct another HTTP request to read the payment.
	reqRead, errRead := http.NewRequest(http.MethodGet, "/v1/payments/"+newID, nil)
	is.NoErr(errRead)

	// Send the read request and assert on the length of the response.
	resp3 := checkResponseStatus(t, s, reqRead, http.StatusOK) // expecting HTTP 200
	defer resp3.Body.Close()

	respBodyBytes, errReadAll := ioutil.ReadAll(resp3.Body)
	is.NoErr(errReadAll)
	is.True(len(string(respBodyBytes)) > 0) // response body should have some content

	// Unmarshal into a slice of Payment structs.
	returnedPayment := server.NewWrapper("")
	errUnmarshal := json.Unmarshal(respBodyBytes, &returnedPayment)
	is.NoErr(errUnmarshal)

	// Assert that the Amount of the Payment returned has been updated appropriately.
	is.Equal(updatedAmount, returnedPayment.Data[0].Attributes.Amount) // amount not updated
}
