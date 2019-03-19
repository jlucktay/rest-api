package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// AlreadyExistsError is returned when a Payment with a given ID exists.
type AlreadyExistsError struct {
	id uuid.UUID
}

func (re *AlreadyExistsError) Error() string {
	return fmt.Sprintf("Payment ID '%s' already exists.", re.id)
}

// NotFoundError is returned when a Payment with a given ID cannot be found.
type NotFoundError struct {
	id uuid.UUID
}

func (re *NotFoundError) Error() string {
	return fmt.Sprintf("Payment ID '%s' not found.", re.id)
}
