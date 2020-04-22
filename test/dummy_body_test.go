package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/matryer/is"

	"go.jlucktay.dev/rest-api/pkg/server"
	"go.jlucktay.dev/rest-api/pkg/storage"
)

func dummyBodyCases(existingID uuid.UUID) map[string]crudTestCase { //nolint:funlen
	return map[string]crudTestCase{
		"Create a new payment with a Payment request body": {
			verb:     http.MethodPost,
			path:     "/v1/payments",
			expected: http.StatusCreated, // 201
		},
		"Create a new payment on a pre-existing ID with a Payment request body": {
			verb:     http.MethodPost,
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			expected: http.StatusConflict, // 409
		},
		"Create a new payment on a non-existent valid ID with a Payment request body": {
			verb:     http.MethodPost,
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			expected: http.StatusNotFound, // 404
		},
		"Create a new payment on an invalid ID with a Payment request body": {
			verb:     http.MethodPost,
			path:     "/v1/payments/not-a-valid-v4-uuid",
			expected: http.StatusNotFound, // 404
		},
		"Update all existing payments": {
			verb:     http.MethodPut,
			path:     "/v1/payments",
			expected: http.StatusMethodNotAllowed, // 405
		},
		"Update an existing payment": {
			verb:     http.MethodPut,
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			expected: http.StatusNoContent, // 204; update is OK, but response has no body/content
		},
		"Update a non-existent payment at a valid ID": {
			verb:     http.MethodPut,
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			expected: http.StatusNotFound, // 404
		},
		"Update a non-existent payment at an invalid ID": {
			verb:     http.MethodPut,
			path:     "/v1/payments/not-a-valid-v4-uuid",
			expected: http.StatusNotFound, // 404
		},
	}
}

// TestDummyBodyCreateUpdate tests creating and updating payment records, with a Payment in the HTTP request body.
func TestDummyBodyCreateUpdate(t *testing.T) {
	existingID := uuid.Must(uuid.NewV4())

	for name, tC := range dummyBodyCases(existingID) {
		// pin! ref: https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		name, tC := name, tC

		is := is.New(t)
		srv := server.New(server.InMemory, false)
		dummyPayment := storage.Payment{Amount: Amount}
		errCreate := srv.Storage.CreateSpecificID(existingID, dummyPayment)
		is.NoErr(errCreate)

		t.Run(name, func(t *testing.T) {
			t.Parallel() // Don't use .Parallel() without pinning.

			var buf bytes.Buffer
			errEncode := json.NewEncoder(&buf).Encode(dummyPayment)
			is.NoErr(errEncode)

			req, err := http.NewRequest(tC.verb, tC.path, &buf)
			is.NoErr(err)

			w := httptest.NewRecorder()
			srv.Router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			is.Equal(resp.StatusCode, tC.expected)
		})
	}
}
