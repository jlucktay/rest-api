package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/matryer/is"

	"go.jlucktay.dev/rest-api/pkg/server"
	"go.jlucktay.dev/rest-api/pkg/storage"
)

func nilBodyCases(existingID uuid.UUID) map[string]crudTestCase { //nolint:funlen
	return map[string]crudTestCase{
		"Create a new payment with an empty request body": {
			verb:     http.MethodPost,
			path:     "/v1/payments",
			expected: http.StatusBadRequest, // 400
		},
		"Create a new payment on a pre-existing ID with an empty request body": {
			verb:     http.MethodPost,
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			expected: http.StatusConflict, // 409
		},
		"Create a new payment on a non-existent valid ID with an empty request body": {
			verb:     http.MethodPost,
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			expected: http.StatusNotFound, // 404
		},
		"Create a new payment on an invalid ID with an empty request body": {
			verb:     http.MethodPost,
			path:     "/v1/payments/not-a-valid-uuid",
			expected: http.StatusNotFound, // 404
		},
		"Read the entire collection of existing payments": {
			verb:     http.MethodGet,
			path:     "/v1/payments",
			expected: http.StatusOK, // 200
		},
		"Read a limited collection of existing payments": {
			verb:     http.MethodGet,
			path:     "/v1/payments?offset=5&limit=5",
			expected: http.StatusOK, // 200
		},
		"Read a single existing payment": {
			verb:     http.MethodGet,
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			expected: http.StatusOK, // 200
		},
		"Read a non-existent payment at a valid ID": {
			verb:     http.MethodGet,
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			expected: http.StatusNotFound, // 404
		},
		"Read a non-existent payment at an invalid ID": {
			verb:     http.MethodGet,
			path:     "/v1/payments/not-a-valid-uuid",
			expected: http.StatusNotFound, // 404
		},
		"Update all existing payments": {
			verb:     http.MethodPut,
			path:     "/v1/payments",
			expected: http.StatusMethodNotAllowed, // 405
		},
		"Update a non-existent payment at an invalid ID": {
			verb:     http.MethodPut,
			path:     "/v1/payments/not-a-valid-uuid",
			expected: http.StatusNotFound, // 404
		},
		"Delete all existing payments": {
			verb:     http.MethodDelete,
			path:     "/v1/payments",
			expected: http.StatusMethodNotAllowed, // 405
		},
		"Delete an existing payment at a valid ID": {
			verb:     http.MethodDelete,
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			expected: http.StatusOK, // 200
		},
		"Delete a non-existent payment at a valid ID": {
			verb:     http.MethodDelete,
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			expected: http.StatusNotFound, // 404
		},
		"Delete a non-existent payment at an invalid ID": {
			verb:     http.MethodDelete,
			path:     "/v1/payments/not-a-valid-uuid",
			expected: http.StatusNotFound, // 404
		},
	}
}

// TestNilBodyCRUD tests creating/reading/updating/deleting payment records, without any body in the HTTP request.
func TestNilBodyCRUD(t *testing.T) {
	existingID := uuid.Must(uuid.NewV4())

	for name, tC := range nilBodyCases(existingID) {
		// pin! ref: https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		name, tC := name, tC

		is := is.New(t)
		srv := server.New(server.InMemory, false)
		existingPayment := storage.Payment{Amount: Amount}
		errCreate := srv.Storage.CreateSpecificID(existingID, existingPayment)
		is.NoErr(errCreate)

		for count := 0; count < 20; count++ {
			_, errCreate := srv.Storage.Create(existingPayment)
			is.NoErr(errCreate)
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel() // Don't use .Parallel() without pinning.

			req, err := http.NewRequest(tC.verb, tC.path, nil)
			is.NoErr(err)

			w := httptest.NewRecorder()
			srv.Router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			is.Equal(resp.StatusCode, tC.expected)
		})
	}
}
