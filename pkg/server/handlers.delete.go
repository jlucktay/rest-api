package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/jlucktay/rest-api/pkg/storage"
)

func (a *Server) deletePaymentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.FromStringOrNil(chi.URLParam(r, "id"))

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
