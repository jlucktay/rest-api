package storage

import (
	uuid "github.com/satori/go.uuid"
)

// PaymentStorage allows storage, retrieval, updating, and deletion of Payment structs.
type PaymentStorage interface {
	Initialise() error
	Terminate() error
	Create(Payment) (uuid.UUID, error)
	CreateSpecificID(uuid.UUID, Payment) error
	Read(uuid.UUID) (Payment, error)
	ReadAll(ReadAllOptions) (map[uuid.UUID]Payment, error)
	Update(uuid.UUID, Payment) error
	Delete(uuid.UUID) error
}

// ReadAllOptions is a config struct for supplying optional parameters to ReadAll.
type ReadAllOptions struct {
	Limit  uint
	Offset uint
}

// DefaultLimit is the default limit on the number of Payments returned by ReadAll.
const DefaultLimit = 10
