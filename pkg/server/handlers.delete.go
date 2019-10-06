package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/jlucktay/rest-api/pkg/storage"
)

func (s *Server) deletePaymentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.FromStringOrNil(chi.URLParam(r, "id"))

		if id == uuid.Nil {
			log.Debug("ID was invalid, returning.")
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		if errDelete := s.Storage.Delete(id); errDelete == nil {
			log.Debugf("Successfully deleted record with given ID '%s'.", id)
			w.WriteHeader(http.StatusOK) // 200
			return
		}

		log.Debug("ID was not found, returning.")
		http.Error(w, (&storage.NotFoundError{ID: id}).Error(), http.StatusNotFound) // 404
	}
}
