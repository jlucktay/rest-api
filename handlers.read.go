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
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		var opts ReadAllOptions

		if limit != "" {
			l, errLimit := strconv.Atoi(limit)
			if errLimit != nil || l <= 0 {
				http.Error(
					w,
					"The 'limit' parameter should be a positive integer.",
					http.StatusBadRequest,
				)
				return
			}
			opts.limit = uint(l)
		}

		if offset != "" {
			o, errOffset := strconv.Atoi(offset)
			if errOffset != nil || o <= 0 {
				http.Error(
					w,
					"The 'offset' parameter should be a positive integer.",
					http.StatusBadRequest,
				)
				return
			}
			opts.offset = uint(o)
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

func (a *apiServer) readPaymentById() httprouter.Handle {
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
