package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/matryer/is"
)

// TestDummyBodyCreateUpdate tests creating and updating payment records, with a simple dummy Payment in the HTTP
// request body.
func TestDummyBodyCreateUpdate(t *testing.T) {
	srv := server.New(server.InMemory)
	existingID := uuid.Must(uuid.NewV4())
	dummyPayment := &storage.Payment{Amount: 123.45}
	errCreate := srv.Storage.CreateSpecificID(existingID, *dummyPayment)
	is.New(t).NoErr(errCreate)

	testCases := []struct {
		desc     string
		path     string
		verb     string
		expected int
	}{
		{
			desc:     "Create a new payment with a Payment request body",
			path:     "/v1/payments",
			verb:     http.MethodPost,
			expected: http.StatusCreated, // 201
		},
		{
			desc:     "Create a new payment on a pre-existing ID with a Payment request body",
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			verb:     http.MethodPost,
			expected: http.StatusConflict, // 409
		},
		{
			desc:     "Create a new payment on a non-existent valid ID with a Payment request body",
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Create a new payment on an invalid ID with a Payment request body",
			path:     "/v1/payments/not-a-valid-v4-uuid",
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Update all existing payments",
			path:     "/v1/payments",
			verb:     http.MethodPut,
			expected: http.StatusMethodNotAllowed, // 405
		},
		{
			desc:     "Update an existing payment",
			path:     fmt.Sprintf("/v1/payments/%s", existingID),
			verb:     http.MethodPut,
			expected: http.StatusNoContent, // 204; update is OK, but response has no body/content
		},
		{
			desc:     "Update a non-existent payment at a valid ID",
			path:     fmt.Sprintf("/v1/payments/%s", uuid.Must(uuid.NewV4())),
			verb:     http.MethodPut,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Update a non-existent payment at an invalid ID",
			path:     "/v1/payments/not-a-valid-v4-uuid",
			verb:     http.MethodPut,
			expected: http.StatusNotFound, // 404
		},
	}
	for _, tC := range testCases {
		tC := tC // pin!
		w := httptest.NewRecorder()

		t.Run(tC.desc, func(t *testing.T) {
			i := is.New(t)

			var buf bytes.Buffer
			errEncode := json.NewEncoder(&buf).Encode(dummyPayment)
			i.NoErr(errEncode)

			req, err := http.NewRequest(tC.verb, tC.path, &buf)
			i.NoErr(err)

			srv.Router.ServeHTTP(w, req)
			i.Equal(tC.expected, w.Result().StatusCode)
		})
	}
}
