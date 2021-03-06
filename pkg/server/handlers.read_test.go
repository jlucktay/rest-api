package server_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/rest-api/pkg/server"
	"go.jlucktay.dev/rest-api/pkg/storage"
	"go.jlucktay.dev/rest-api/test"
)

func TestReadMultiplePayments(t *testing.T) {
	var w *httptest.ResponseRecorder

	s := server.New(server.InMemory, false)
	is := is.New(t)

	// Construct a HTTP request which creates a payment.
	p := storage.Payment{Amount: test.Amount}
	j, errMarshal := json.Marshal(p)
	is.NoErr(errMarshal)

	// Send it multiple times, to create multiple payments.
	for count := 0; count < 250; count++ {
		reqCreate, errCreate := http.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBuffer(j))
		is.NoErr(errCreate)
		reqCreate.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		s.Router.ServeHTTP(w, reqCreate)

		resp := w.Result()
		defer resp.Body.Close()

		is.Equal(http.StatusCreated, resp.StatusCode) // expecting HTTP 201
	}

	// Construct another HTTP request to read the payments.
	// 'limit' and 'offset' are not specified here, so we are falling back onto the default values.
	reqRead, errRead := http.NewRequest(http.MethodGet, "/v1/payments", nil)
	is.NoErr(errRead)

	// Read the payments.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqRead)

	resp := w.Result()
	defer resp.Body.Close()

	respBodyBytes, errReadAll := ioutil.ReadAll(resp.Body)
	is.NoErr(errReadAll)
	is.True(len(string(respBodyBytes)) > 0) // response body should have some content

	// Unmarshal them into a slice of Payment structs.
	returnedPayments := server.NewWrapper("")
	errUnmarshal := json.Unmarshal(respBodyBytes, &returnedPayments)
	is.NoErr(errUnmarshal)

	// Assert on the number of Payment structs returned.
	is.Equal(storage.DefaultLimit, len(returnedPayments.Data)) // check default pagination limit
}
