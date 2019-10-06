package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/matryer/is"
)

func TestCreateNewPayment(t *testing.T) {
	s := server.New(server.InMemory, false)
	var w *httptest.ResponseRecorder
	i := is.New(t)

	// Construct a HTTP request which creates a payment.
	p := storage.Payment{Amount: 123.45}
	j, errMarshal := json.Marshal(p)
	i.NoErr(errMarshal)
	reqBody := bytes.NewBuffer(j)
	reqCreate, errCreate := http.NewRequest(http.MethodPost, "/v1/payments", reqBody)
	i.NoErr(errCreate)
	reqCreate.Header.Set("Content-Type", "application/json")

	// Send it, and record the HTTP back and forth.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqCreate)
	resp := w.Result()
	defer resp.Body.Close()
	i.Equal(http.StatusCreated, resp.StatusCode)

	// Make sure the response had a Location header pointing at the new payment.
	loc := resp.Header.Get("Location")
	r := regexp.MustCompile("^/v1/payments/([0-9a-f-]{36})$")
	i.True(r.MatchString(loc))
	newID := r.FindStringSubmatch(loc)[1]

	// Construct another HTTP request to read the new payment.
	reqRead, errRead := http.NewRequest(http.MethodGet, "/v1/payments/"+newID, nil)
	i.NoErr(errRead)

	// Read the new payment using the ID returned.
	w = httptest.NewRecorder()
	s.Router.ServeHTTP(w, reqRead)
	resp2 := w.Result()
	defer resp2.Body.Close()
	i.Equal(http.StatusOK, resp2.StatusCode)
}
