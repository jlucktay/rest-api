package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

func (a *apiServer) readPayments() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var opts ReadAllOptions
		if errLimit := applyFromQuery(r, "limit", &opts.limit); errLimit != nil {
			http.Error(w, errLimit.Error(), http.StatusBadRequest)
			return
		}
		if errOffset := applyFromQuery(r, "offset", &opts.offset); errOffset != nil {
			http.Error(w, errOffset.Error(), http.StatusBadRequest)
			return
		}

		allPayments, errRead := a.storage.ReadAll(opts)
		if errRead != nil {
			http.Error(
				w,
				fmt.Sprintf("Error reading all payments: %s", errRead),
				http.StatusInternalServerError, // 500
			)
			return
		}

		var wrappedPayments readWrapper
		wrappedPayments.init(r)
		for id, payment := range allPayments {
			wrappedPayments.addPayment(id, payment)
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

func (a *apiServer) readPaymentByID() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := uuid.FromStringOrNil(p.ByName("id"))

		if id == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		if payRead, errRead := a.storage.Read(id); errRead == nil {
			var wrappedPayment readWrapper
			wrappedPayment.init(r)
			wrappedPayment.addPayment(id, payRead)

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

		http.Error(w, (&NotFoundError{id}).Error(), http.StatusNotFound) // 404
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
