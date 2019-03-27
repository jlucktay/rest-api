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

func TestDocumentationSingle(t *testing.T) {
	// set up server
	srv := newAPIServer(InMemory)
	w := httptest.NewRecorder()
	i := is.New(t)

	// put the documentation payment into it
	singleBytes, errReadFile := ioutil.ReadFile("testdata/doc.single.json")
	i.NoErr(errReadFile)
	var single Payment
	errUmSingle := json.Unmarshal(singleBytes, &single)
	i.NoErr(errUmSingle)
	existingID := uuid.Must(uuid.NewV4())
	errCreate := srv.storage.createSpecificID(existingID, single)
	i.NoErr(errCreate)

	// do a http request for the ID of the payment
	req, errReq := http.NewRequest(http.MethodGet, fmt.Sprintf("/payments/%s", existingID), nil)
	i.NoErr(errReq)
	srv.router.ServeHTTP(w, req)
	i.Equal(http.StatusOK, w.Result().StatusCode)

	// assert that it matches the documentation JSON
	var expected readWrapper
	expected.init(req)
	expected.addPayment(existingID, single)
	responseBytes, errReadResponse := ioutil.ReadAll(w.Result().Body)
	i.NoErr(errReadResponse)

	fmt.Printf("responseBytes: '%+v'\n", string(responseBytes))

	var actual readWrapper
	// actual.init(req)
	errUmResponse := json.Unmarshal(responseBytes, &actual.Data)
	i.NoErr(errUmResponse)

	fmt.Printf("expected: '%+v'\n", expected)
	fmt.Printf("actual: '%+v'\n", actual)

	i.True(reflect.DeepEqual(expected, actual))
}
