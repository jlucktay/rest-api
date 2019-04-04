package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/matryer/is"
	uuid "github.com/satori/go.uuid"
)

type testDataWrapper struct {
	Attributes storage.Payment `json:"attributes"`
	ID         uuid.UUID       `json:"id"`
}

func TestDocumentation(t *testing.T) {
	testCases := []struct {
		desc         string
		testdataFile string
		getPath      string
	}{
		{
			desc:         "Get a single payment.",
			testdataFile: "testdata/doc.single.json",
			getPath:      "/payments/97fe60ba-1334-439f-91db-32cc3cde036a",
		},
		{
			desc:         "Get multiple payments.",
			testdataFile: "testdata/doc.multiple.json",
			getPath:      "/payments",
		},
	}

	for _, tC := range testCases {
		tC := tC // pin!
		t.Run(tC.desc, func(t *testing.T) {
			// Set up an API server to test against.
			srv := server.New(server.InMemory)
			w := httptest.NewRecorder()
			i := is.New(t)

			// POST the payment(s) from the testdata JSON file into the server.
			fileBytes, errReadFile := ioutil.ReadFile(tC.testdataFile)
			i.NoErr(errReadFile)
			var wrapped []testDataWrapper
			errUm := json.Unmarshal(fileBytes, &wrapped)
			i.NoErr(errUm)
			for _, testdata := range wrapped {
				errCreate := srv.Storage.CreateSpecificID(testdata.ID, testdata.Attributes)
				i.NoErr(errCreate)
			}

			// Do a HTTP request to GET the payment(s).
			req, errReq := http.NewRequest(http.MethodGet, tC.getPath, nil)
			i.NoErr(errReq)
			srv.Router.ServeHTTP(w, req)
			i.Equal(http.StatusOK, w.Result().StatusCode)

			// Put info from the testdata JSON file directly into a wrapper struct.
			expected := server.NewWrapper(req.URL.String())
			for _, testdata := range wrapped {
				expected.AddPayment(testdata.ID, testdata.Attributes)
			}

			// Assert that it matches the JSON returned by the API.
			responseBytes, errReadResponse := ioutil.ReadAll(w.Result().Body)
			i.NoErr(errReadResponse)
			actual := server.NewWrapper(req.URL.String())
			errUmResponse := json.Unmarshal(responseBytes, &actual)
			i.NoErr(errUmResponse)
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Fatalf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
