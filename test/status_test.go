package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/matryer/is"
)

func TestStatusCode(t *testing.T) {
	// Arrange
	testCases := []struct {
		desc     string
		path     string
		verb     string
		expected int
	}{
		{
			desc:     "Create a new payment on a non-existent valid ID",
			path:     "/payments/60c4feb1-bf67-488a-8d04-627bac487c05",
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Create a new payment on an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Read the entire collection of existing payments",
			path:     "/payments",
			verb:     http.MethodGet,
			expected: http.StatusOK, // 200
		},
		{
			desc:     "Read a non-existent payment at a valid ID",
			path:     "/payments/29e1c453-8cc7-47b8-9c48-7e44b4f9ba26",
			verb:     http.MethodGet,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Read a non-existent payment at an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodGet,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Update all existing payments",
			path:     "/payments",
			verb:     http.MethodPut,
			expected: http.StatusMethodNotAllowed, // 405
		},
		{
			desc:     "Update a non-existent payment at an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodPut,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Delete all existing payments",
			path:     "/payments",
			verb:     http.MethodDelete,
			expected: http.StatusMethodNotAllowed, // 405
		},
		{
			desc:     "Delete a non-existent payment at a valid ID",
			path:     "/payments/943c4811-f66a-4cb1-8d5b-3ed7da0ce934",
			verb:     http.MethodDelete,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Delete a non-existent payment at an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodDelete,
			expected: http.StatusNotFound, // 404
		},
	}

	srv := server.New(server.InMemory)

	// Act & Assert
	for _, tC := range testCases {
		tC := tC // pin!
		w := httptest.NewRecorder()

		t.Run(tC.desc, func(t *testing.T) {
			i := is.New(t)

			req, err := http.NewRequest(tC.verb, tC.path, nil)
			i.NoErr(err)

			srv.Router.ServeHTTP(w, req)
			i.Equal(tC.expected, w.Result().StatusCode)
		})
	}
}