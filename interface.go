package main

import (
	uuid "github.com/satori/go.uuid"
)

// PaymentStorage allows storage, retrieval, updating, and deletion of Payment
// structs.
type PaymentStorage interface {
	Init() error
	Create(Payment) (uuid.UUID, error)
	createSpecificId(id uuid.UUID, p Payment) error
	Read(uuid.UUID) (Payment, error)
	ReadAll() ([]Payment, error)
	Update(uuid.UUID, Payment) error
	Delete(uuid.UUID) error
}
