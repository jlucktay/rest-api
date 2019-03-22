package main

import (
	"encoding/json"
	"errors"
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
		if errLimit := applyFromQuery(r.URL.Query().Get("limit"), &opts.limit); errLimit != nil {
			http.Error(w, "'limit':"+errLimit.Error(), http.StatusBadRequest)
			return
		}
		if errOffset := applyFromQuery(r.URL.Query().Get("offset"), &opts.offset); errOffset != nil {
			http.Error(w, "'offset':"+errOffset.Error(), http.StatusBadRequest)
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

		allBytes, errMarshal := json.Marshal(allPayments)
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
			payBytes, errMarshal := json.Marshal(payRead)
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
func applyFromQuery(input string, setting *uint) error {
	if input != "" {
		i, errConvert := strconv.Atoi(input)
		if errConvert != nil || i <= 0 {
			return errors.New("The query parameter should be a positive integer.")
		}
		*setting = uint(i)
	}
	return nil
}
