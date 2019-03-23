package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

func (a *apiServer) deletePaymentByID() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := uuid.FromStringOrNil(p.ByName("id"))

		if id == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		if errDelete := a.storage.Delete(id); errDelete == nil {
			w.WriteHeader(http.StatusOK) // 200
			return
		}

		http.Error(w, (&NotFoundError{id}).Error(), http.StatusNotFound) // 404
	}
}
