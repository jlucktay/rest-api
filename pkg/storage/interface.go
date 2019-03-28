package storage

import (
	uuid "github.com/satori/go.uuid"
)

// PaymentStorage allows storage, retrieval, updating, and deletion of Payment
// structs.
type PaymentStorage interface {
	Init() error
	Create(Payment) (uuid.UUID, error)
	createSpecificID(uuid.UUID, Payment) error
	Read(uuid.UUID) (Payment, error)
	ReadAll(ReadAllOptions) (map[uuid.UUID]Payment, error)
	Update(uuid.UUID, Payment) error
	Delete(uuid.UUID) error
}

// ReadAllOptions is a config struct for calling 'PaymentStorage.ReadAll()'.
type ReadAllOptions struct {
	Limit  uint
	Offset uint
}
