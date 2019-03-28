package storage

import (
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router  *httprouter.Router
	storage storage.PaymentStorage
}

// StorageType is an enum to differentiate between storage implementations.
type StorageType byte

const (
	// InMemory describes a storage system that is held in memory only, and not
	// persisted to disk.
	InMemory StorageType = iota

	// Mongo describes a storage system that persists Payment records in a
	// MongoDB database.
	Mongo
)
