package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

func (s *Server) createPayments() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.ContentLength == 0 {
			http.Error(w, "Empty request body.", http.StatusBadRequest) // 400
			return
		}

		bodyBytes, errRead := ioutil.ReadAll(r.Body)
		if errRead != nil {
			log.Fatal(errRead)
		}
		defer r.Body.Close()

		var p storage.Payment
		errUm := json.Unmarshal(bodyBytes, &p)
		if errUm != nil {
			log.Fatal(errUm)
		}

		id, errCreate := s.Storage.Create(p)
		if errCreate != nil {
			log.Fatal(errCreate)
		}

		w.Header().Set("Location", fmt.Sprintf("/payments/%s", id))
		w.WriteHeader(http.StatusCreated) // 201
	}
}

func (s *Server) createPaymentByID() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := uuid.FromStringOrNil(p.ByName("id"))

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