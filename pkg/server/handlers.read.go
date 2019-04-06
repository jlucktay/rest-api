package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) readPayments() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
			http.Error(
				w,
				fmt.Sprintf("Error reading all payments: %s", errRead),
				http.StatusInternalServerError, // 500
			)
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
			http.Error(
				w,
				fmt.Sprintf("Error marshaling payments: %s", errMarshal),
				http.StatusInternalServerError, // 500
			)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		w.Header().Set("Content-Type", "application/json")
		if _, errWrite := w.Write(allBytes); errWrite != nil {
			log.Fatal(errWrite)
		}
	}
}

func (s *Server) readPaymentByID() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := uuid.FromStringOrNil(p.ByName("id"))

		if id == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		if payRead, errRead := s.Storage.Read(id); errRead == nil {
			wrappedPayment := NewWrapper(r.URL.String())
			wrappedPayment.AddPayment(id, payRead)

			payBytes, errMarshal := json.Marshal(wrappedPayment)
			if errMarshal != nil {
				log.Fatal(errMarshal)
			}

			w.WriteHeader(http.StatusOK) // 200
			w.Header().Set("Content-Type", "application/json")
			if _, errWrite := w.Write(payBytes); errWrite != nil {
				log.Fatal(errWrite)
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
