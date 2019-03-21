package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// AlreadyExistsError is returned when a Payment with a given ID exists.
type AlreadyExistsError struct {
	id uuid.UUID
}

func (aee *AlreadyExistsError) Error() string {
	return fmt.Sprintf("Payment ID '%s' already exists.", aee.id)
}

// NotFoundError is returned when a Payment with a given ID cannot be found.
type NotFoundError struct {
	id uuid.UUID
}

func (nfe *NotFoundError) Error() string {
	return fmt.Sprintf("Payment ID '%s' not found.", nfe.id)
}

// OffsetOutOfBounds is returned when the 'offset' parameter inside a
// ReadAllOptions struct exceeds the number of elements available in
// PaymentStorage.
type OffsetOutOfBounds struct {
	offset uint
}

func (ooob *OffsetOutOfBounds) Error() string {
	return fmt.Sprintf("Offset '%d' is out of bounds.", ooob.offset)
}
