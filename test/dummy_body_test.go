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

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
)

// TestDummyBodyCreateUpdate tests creating and updating payment records, with a simple dummy Payment in the HTTP
// request body.
func TestDummyBodyCreateUpdate(t *testing.T) {
	existingID := uuid.Must(uuid.NewV4())

	testCases := map[string]struct {
		path     string
		verb     string
		expected int
	}{
		"Create a new payment with a Payment request body": {
			path:     "/v1/payments",
			verb:     http.MethodPost,
			expected: http.StatusCreated, // 201
		},
		"Create a new payment on a pre-existing ID with a Payment request body": {
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			verb:     http.MethodPost,
			expected: http.StatusConflict, // 409
		},
		"Create a new payment on a non-existent valid ID with a Payment request body": {
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		"Create a new payment on an invalid ID with a Payment request body": {
			path:     "/v1/payments/not-a-valid-v4-uuid",
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		"Update all existing payments": {
			path:     "/v1/payments",
			verb:     http.MethodPut,
			expected: http.StatusMethodNotAllowed, // 405
		},
		"Update an existing payment": {
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			verb:     http.MethodPut,
			expected: http.StatusNoContent, // 204; update is OK, but response has no body/content
		},
		"Update a non-existent payment at a valid ID": {
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			verb:     http.MethodPut,
			expected: http.StatusNotFound, // 404
		},
		"Update a non-existent payment at an invalid ID": {
			path:     "/v1/payments/not-a-valid-v4-uuid",
			verb:     http.MethodPut,
			expected: http.StatusNotFound, // 404
		},
	}
	for name, tC := range testCases {
		tC := tC // pin!

		srv := server.New(server.InMemory, false)
		dummyPayment := &storage.Payment{Amount: Amount}
		errCreate := srv.Storage.CreateSpecificID(existingID, *dummyPayment)
		is.New(t).NoErr(errCreate)

		w := httptest.NewRecorder()

		t.Run(name, func(t *testing.T) {
			is := is.New(t)

			var buf bytes.Buffer
			errEncode := json.NewEncoder(&buf).Encode(dummyPayment)
			is.NoErr(errEncode)

			req, err := http.NewRequest(tC.verb, tC.path, &buf)
			is.NoErr(err)

			srv.Router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			is.Equal(resp.StatusCode, tC.expected)
		})
	}
}
