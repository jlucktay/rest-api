package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"

	"github.com/jlucktay/rest-api/pkg/storage"
)

func (s *Server) updatePaymentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.FromStringOrNil(chi.URLParam(r, "id"))

		if id == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		if r.ContentLength == 0 {
			http.Error(w, "Empty request body.", http.StatusBadRequest) // 400
			return
		}

		bodyBytes, errRead := ioutil.ReadAll(r.Body)
		if errRead != nil {
			logrus.Fatal(errRead)
		}
		defer r.Body.Close()

		var payment storage.Payment
		errUm := json.Unmarshal(bodyBytes, &payment)
		if errUm != nil {
			logrus.Fatal(errUm)
		}

		if errUpdate := s.Storage.Update(id, payment); errUpdate == nil {
			w.WriteHeader(http.StatusNoContent) // 204
			return
		}

		http.Error(w, (&storage.NotFoundError{ID: id}).Error(), http.StatusNotFound) // 404
	}
}
