package server

import (
	"net/http"

	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

func (a *Server) deletePaymentByID() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := uuid.FromStringOrNil(p.ByName("id"))

		if id == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		if errDelete := a.Storage.Delete(id); errDelete == nil {
			w.WriteHeader(http.StatusOK) // 200
			return
		}

		http.Error(w, (&storage.NotFoundError{ID: id}).Error(), http.StatusNotFound) // 404
	}
}
