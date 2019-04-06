package server

import (
	"github.com/go-chi/chi"
	"github.com/jlucktay/rest-api/pkg/storage"
)

// Server is a RESTful HTTP API server offering CRUD functionality to store Payments.
type Server struct {
	Router  *chi.Mux
	Storage storage.PaymentStorage
}

// StorageType is an enum to differentiate between storage implementations.
type StorageType byte

const (
	// InMemory describes a storage system that is held in memory only, and not persisted to disk.
	InMemory StorageType = iota

	// Mongo describes a storage system that persists Payment records in a MongoDB database.
	Mongo
)

const (
	defaultDatabase   = "rest-api"
	defaultCollection = "payments"
)
