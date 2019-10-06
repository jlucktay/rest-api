package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"

	"github.com/jlucktay/rest-api/pkg/storage"
)

func (s *Server) readPayments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var opts storage.ReadAllOptions
		if errLimit := applyFromQuery(r, "limit", &opts.Limit); errLimit != nil {
			http.Error(w, errLimit.Error(), http.StatusBadRequest)
			return
		}
		if errOffset := applyFromQuery(r, "offset", &opts.Offset); errOffset != nil {
			http.Error(w, errOffset.Error(), http.StatusBadRequest)
			return
		}

		allPayments, errRead := s.Storage.ReadAll(opts)
		if errRead != nil {
			http.Error(w, fmt.Sprintf("Error reading all: %s", errRead), http.StatusInternalServerError) // 500
			return
		}

		keys := []string{}
		for a := range allPayments {
			keys = append(keys, a.String())
		}
		sort.Strings(keys)

		wrappedPayments := NewWrapper(r.URL.String())
		for _, sID := range keys {
			id := uuid.FromStringOrNil(sID)
			wrappedPayments.AddPayment(id, allPayments[id])
		}

		allBytes, errMarshal := json.Marshal(wrappedPayments)
		if errMarshal != nil {
			http.Error(w, fmt.Sprintf("Error marshaling: %s", errMarshal), http.StatusInternalServerError) // 500
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		w.Header().Set("Content-Type", "application/json")
		if _, errWrite := w.Write(allBytes); errWrite != nil {
			logrus.Fatal(errWrite)
		}
	}
}

func (s *Server) readPaymentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.FromStringOrNil(chi.URLParam(r, "id"))

		if id == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		if payRead, errRead := s.Storage.Read(id); errRead == nil {
			wrappedPayment := NewWrapper(r.URL.String())
			wrappedPayment.AddPayment(id, payRead)

			payBytes, errMarshal := json.Marshal(wrappedPayment)
			if errMarshal != nil {
				logrus.Fatal(errMarshal)
			}

			w.WriteHeader(http.StatusOK) // 200
			w.Header().Set("Content-Type", "application/json")
			if _, errWrite := w.Write(payBytes); errWrite != nil {
				logrus.Fatal(errWrite)
			}
			return
		}

		http.Error(w, (&storage.NotFoundError{ID: id}).Error(), http.StatusNotFound) // 404
	}
}

// applyFromQuery takes a string (from an HTTP request query) as well as a
// pointer to a uint which it will apply the value of the string to.
func applyFromQuery(req *http.Request, query string, setting *uint) error {
	if q := req.URL.Query().Get(query); q != "" {
		i, errConvert := strconv.Atoi(q)
		if errConvert != nil || i <= 0 {
			return fmt.Errorf("the query parameter '%s' should be a positive integer", query)
		}
		*setting = uint(i)
	}
	return nil
}
