package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/matryer/is"
	uuid "github.com/satori/go.uuid"
)

type testDataWrapper struct {
	Attributes Payment   `json:"attributes"`
	ID         uuid.UUID `json:"id"`
}

func TestDocumentationSingle(t *testing.T) {
	// Set up an API server to test against.
	srv := newAPIServer(InMemory)
	w := httptest.NewRecorder()
	i := is.New(t)

	// Put the single payment from the documentation into the server.
	singleBytes, errReadFile := ioutil.ReadFile("testdata/doc.single.json")
	i.NoErr(errReadFile)
	var single testDataWrapper
	errUmSingle := json.Unmarshal(singleBytes, &single)
	i.NoErr(errUmSingle)
	errCreate := srv.storage.createSpecificID(single.ID, single.Attributes)
	i.NoErr(errCreate)

	// Do a HTTP request for the single payment.
	req, errReq := http.NewRequest(http.MethodGet, fmt.Sprintf("/payments/%s", single.ID), nil)
	i.NoErr(errReq)
	srv.router.ServeHTTP(w, req)
	i.Equal(http.StatusOK, w.Result().StatusCode)

	// Put info from the ./testdata/ JSON file into a wrapper struct.
	var expected readWrapper
	expected.init(req)
	expected.addPayment(single.ID, single.Attributes)

	// Assert that it matches the JSON returned by the API.
	responseBytes, errReadResponse := ioutil.ReadAll(w.Result().Body)
	i.NoErr(errReadResponse)
	var actual readWrapper
	errUmResponse := json.Unmarshal(responseBytes, &actual)
	i.NoErr(errUmResponse)
	i.True(reflect.DeepEqual(expected, actual))
}
