package main

import (
	"github.com/julienschmidt/httprouter"
)

type apiServer struct {
	router  *httprouter.Router
	storage PaymentStorage
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
