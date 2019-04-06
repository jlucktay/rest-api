package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) updatePaymentByID() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := uuid.FromStringOrNil(p.ByName("id"))

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
			log.Fatal(errRead)
		}
		defer r.Body.Close()

		var payment storage.Payment
		errUm := json.Unmarshal(bodyBytes, &payment)
		if errUm != nil {
			log.Fatal(errUm)
		}

		if errUpdate := s.Storage.Update(id, payment); errUpdate == nil {
			w.WriteHeader(http.StatusNoContent) // 204
			return
		}

		http.Error(w, (&storage.NotFoundError{ID: id}).Error(), http.StatusNotFound) // 404
	}
}
