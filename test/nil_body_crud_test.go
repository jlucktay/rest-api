package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/matryer/is"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
)

// TestNilBodyCRUD tests creating, reading, updating, and deleting payment records, without any body in the HTTP
// requests.
//nolint:funlen // TODO
func TestNilBodyCRUD(t *testing.T) {
	existingID := uuid.Must(uuid.NewV4())

	testCases := map[string]struct {
		path     string
		verb     string
		expected int
	}{
		"Create a new payment with an empty request body": {
			path:     "/v1/payments",
			verb:     http.MethodPost,
			expected: http.StatusBadRequest, // 400
		},
		"Create a new payment on a pre-existing ID with an empty request body": {
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			verb:     http.MethodPost,
			expected: http.StatusConflict, // 409
		},
		"Create a new payment on a non-existent valid ID with an empty request body": {
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		"Create a new payment on an invalid ID with an empty request body": {
			path:     "/v1/payments/not-a-valid-uuid",
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		"Read the entire collection of existing payments": {
			path:     "/v1/payments",
			verb:     http.MethodGet,
			expected: http.StatusOK, // 200
		},
		"Read a limited collection of existing payments": {
			path:     "/v1/payments?offset=5&limit=5",
			verb:     http.MethodGet,
			expected: http.StatusOK, // 200
		},
		"Read a single existing payment": {
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			verb:     http.MethodGet,
			expected: http.StatusOK, // 200
		},
		"Read a non-existent payment at a valid ID": {
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			verb:     http.MethodGet,
			expected: http.StatusNotFound, // 404
		},
		"Read a non-existent payment at an invalid ID": {
			path:     "/v1/payments/not-a-valid-uuid",
			verb:     http.MethodGet,
			expected: http.StatusNotFound, // 404
		},
		"Update all existing payments": {
			path:     "/v1/payments",
			verb:     http.MethodPut,
			expected: http.StatusMethodNotAllowed, // 405
		},
		"Update a non-existent payment at an invalid ID": {
			path:     "/v1/payments/not-a-valid-uuid",
			verb:     http.MethodPut,
			expected: http.StatusNotFound, // 404
		},
		"Delete all existing payments": {
			path:     "/v1/payments",
			verb:     http.MethodDelete,
			expected: http.StatusMethodNotAllowed, // 405
		},
		"Delete an existing payment at a valid ID": {
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			verb:     http.MethodDelete,
			expected: http.StatusOK, // 200
		},
		"Delete a non-existent payment at a valid ID": {
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			verb:     http.MethodDelete,
			expected: http.StatusNotFound, // 404
		},
		"Delete a non-existent payment at an invalid ID": {
			path:     "/v1/payments/not-a-valid-uuid",
			verb:     http.MethodDelete,
			expected: http.StatusNotFound, // 404
		},
	}
	for name, tC := range testCases {
		tC := tC // pin!

		srv := server.New(server.InMemory, false)
		existingPayment := storage.Payment{Amount: Amount}
		errCreate := srv.Storage.CreateSpecificID(existingID, existingPayment)
		is.New(t).NoErr(errCreate)

		for count := 0; count < 20; count++ {
			_, errCreate := srv.Storage.Create(existingPayment)
			is.New(t).NoErr(errCreate)
		}

		w := httptest.NewRecorder()

		t.Run(name, func(t *testing.T) {
			is := is.New(t)

			req, err := http.NewRequest(tC.verb, tC.path, nil)
			is.NoErr(err)

			srv.Router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			is.Equal(resp.StatusCode, tC.expected)
		})
	}
}
