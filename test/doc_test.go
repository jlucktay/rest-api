// Package test runs test coverage over the API router from the outside looking in.
package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
)

type testDataWrapper struct {
	Attributes storage.Payment `json:"attributes"`
	ID         uuid.UUID       `json:"id"`
}

// TestDocumentation is the implementation of the examples from the documentation.
func TestDocumentation(t *testing.T) {
	testCases := map[string]struct {
		testdataFile string
		getPath      string
	}{
		"Get a single payment.": {
			testdataFile: "testdata/doc.single.json",
			getPath:      "/v1/payments/97fe60ba-1334-439f-91db-32cc3cde036a",
		},
		"Get multiple payments.": {
			testdataFile: "testdata/doc.multiple.json",
			getPath:      "/v1/payments",
		},
	}

	for name, tC := range testCases {
		tC := tC // pin!

		t.Run(name, func(t *testing.T) {
			// Set up an API server to test against.
			srv := server.New(server.InMemory, false)
			w := httptest.NewRecorder()
			is := is.New(t)

			// POST the payment(s) from the testdata JSON file into the server.
			fileBytes, errReadFile := ioutil.ReadFile(tC.testdataFile)
			is.NoErr(errReadFile)
			var wrapped []testDataWrapper
			errUm := json.Unmarshal(fileBytes, &wrapped)
			is.NoErr(errUm)
			for _, testdata := range wrapped {
				errCreate := srv.Storage.CreateSpecificID(testdata.ID, testdata.Attributes)
				is.NoErr(errCreate)
			}

			// Do a HTTP request to GET the payment(s).
			req, errReq := http.NewRequest(http.MethodGet, tC.getPath, nil)
			is.NoErr(errReq)
			srv.Router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			is.Equal(resp.StatusCode, http.StatusOK) // expecting HTTP 200

			// Put info from the testdata JSON file directly into a wrapper struct.
			expected := server.NewWrapper(req.URL.String())
			for _, testdata := range wrapped {
				expected.AddPayment(testdata.ID, testdata.Attributes)
			}

			// Assert that it matches the JSON returned by the API.
			responseBytes, errReadResponse := ioutil.ReadAll(resp.Body)
			is.NoErr(errReadResponse)
			actual := server.NewWrapper(req.URL.String())
			errUmResponse := json.Unmarshal(responseBytes, &actual)
			is.NoErr(errUmResponse)
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Fatalf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
