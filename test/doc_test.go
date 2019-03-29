package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/matryer/is"
	uuid "github.com/satori/go.uuid"
)

type testDataWrapper struct {
	Attributes storage.Payment `json:"attributes"`
	ID         uuid.UUID       `json:"id"`
}

func TestDocumentationSingle(t *testing.T) {
	// Set up an API server to test against.
	srv := server.New(server.InMemory)
	w := httptest.NewRecorder()
	i := is.New(t)

	// Put the single payment from the documentation into the server.
	singleBytes, errReadFile := ioutil.ReadFile("testdata/doc.single.json")
	i.NoErr(errReadFile)
	var single testDataWrapper
	errUmSingle := json.Unmarshal(singleBytes, &single)
	i.NoErr(errUmSingle)
	errCreate := srv.Storage.CreateSpecificID(single.ID, single.Attributes)
	i.NoErr(errCreate)

	// Do a HTTP request for the single payment.
	req, errReq := http.NewRequest(http.MethodGet, fmt.Sprintf("/payments/%s", single.ID), nil)
	i.NoErr(errReq)
	srv.Router.ServeHTTP(w, req)
	i.Equal(http.StatusOK, w.Result().StatusCode)

	// Put info from the ./testdata/ JSON file into a wrapper struct.
	var expected server.ReadWrapper
	expected.Init(req)
	expected.AddPayment(single.ID, single.Attributes)

	// Assert that it matches the JSON returned by the API.
	responseBytes, errReadResponse := ioutil.ReadAll(w.Result().Body)
	i.NoErr(errReadResponse)
	var actual server.ReadWrapper
	errUmResponse := json.Unmarshal(responseBytes, &actual)
	i.NoErr(errUmResponse)
	i.True(reflect.DeepEqual(expected, actual))
}

func TestDocumentationMultiple(t *testing.T) {
	// Set up an API server to test against.
	srv := server.New(server.InMemory)
	w := httptest.NewRecorder()
	i := is.New(t)

	// Put the multiple payments from the documentation into the server.
	multipleBytes, errReadFile := ioutil.ReadFile("testdata/doc.multiple.json")
	i.NoErr(errReadFile)
	var multiple []testDataWrapper
	errUmMultiple := json.Unmarshal(multipleBytes, &multiple)
	i.NoErr(errUmMultiple)
	for _, testdata := range multiple {
		errCreate := srv.Storage.CreateSpecificID(testdata.ID, testdata.Attributes)
		i.NoErr(errCreate)
	}

	// Do a HTTP request for the multiple payments.
	req, errReq := http.NewRequest(http.MethodGet, "/payments", nil)
	i.NoErr(errReq)
	srv.Router.ServeHTTP(w, req)
	i.Equal(http.StatusOK, w.Result().StatusCode)

	// Put info from the ./testdata/ JSON file into a wrapper struct.
	var expected server.ReadWrapper
	expected.Init(req)
	for _, testdata := range multiple {
		expected.AddPayment(testdata.ID, testdata.Attributes)
	}

	// Assert that it matches the JSON returned by the API.
	responseBytes, errReadResponse := ioutil.ReadAll(w.Result().Body)
	i.NoErr(errReadResponse)
	var actual server.ReadWrapper
	errUmResponse := json.Unmarshal(responseBytes, &actual)
	i.NoErr(errUmResponse)
	i.True(reflect.DeepEqual(expected, actual))
}
