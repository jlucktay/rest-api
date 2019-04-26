package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/sirupsen/logrus"
)

func (s *Server) createPayments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			http.Error(w, "Empty request body.", http.StatusBadRequest) // 400
			return
		}

		bodyBytes, errRead := ioutil.ReadAll(r.Body)
		if errRead != nil {
			logrus.Fatal(errRead)
		}
		defer r.Body.Close()

		var p storage.Payment
		errUm := json.Unmarshal(bodyBytes, &p)
		if errUm != nil {
			logrus.Fatal(errUm)
		}

		id, errCreate := s.Storage.Create(p)
		if errCreate != nil {
			logrus.Fatal(errCreate)
		}

		w.Header().Set("Location", fmt.Sprintf("/v1/payments/%s", id))
		w.WriteHeader(http.StatusCreated) // 201
	}
}

func (s *Server) createPaymentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.FromStringOrNil(chi.URLParam(r, "id"))

		if id == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		_, errRead := s.Storage.Read(id)
		if errRead == nil {
			http.Error(w, (&storage.AlreadyExistsError{ID: id}).Error(), http.StatusConflict) // 409
			return
		}

		http.Error(w, `Cannot specify an ID for payment creation.
One will be generated for you.`, http.StatusNotFound) // 404
	}
}
