package server_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/matryer/is"
	"github.com/shopspring/decimal"
)

func TestReadMultiplePayments(t *testing.T) {
	s := server.New(server.InMemory)
	var w *httptest.ResponseRecorder
	i := is.New(t)

	// Construct a HTTP request which creates a payment.
	p := storage.Payment{Amount: decimal.NewFromFloat(123.45)}
	j, errMarshal := json.Marshal(p)
	i.NoErr(errMarshal)

	// Send it multiple times, to create multiple payments.
	for count := 0; count < 250; count++ {
		reqCreate, errCreate := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(j))
		i.NoErr(errCreate)
		reqCreate.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		s.Router.ServeHTTP(w, reqCreate)
		i.Equal(http.StatusCreated, w.Result().StatusCode)
	}

	// Construct another HTTP request to read the payments.
	// 'limit' and 'offset' are not specified here, so we are falling back onto
	// the default values.
	reqRead, errRead := http.NewRequest(http.MethodGet, "/payments", nil)
	i.NoErr(errRead)

	// Read the payments.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqRead)
	respBodyBytes, errReadAll := ioutil.ReadAll(w.Result().Body)
	i.NoErr(errReadAll)
	i.True(len(string(respBodyBytes)) > 0)

	// Unmarshal them into a slice of Payment structs.
	returnedPayments := server.NewWrapper("")
	errUnmarshal := json.Unmarshal(respBodyBytes, &returnedPayments)
	i.NoErr(errUnmarshal)

	// Assert on the number of Payment structs returned.
	i.Equal(storage.DefaultLimit, len(returnedPayments.Data))
}
