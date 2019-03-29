// Package org represents an organisation that makes payment transactions.
package org

import (
	uuid "github.com/satori/go.uuid"
)

// ID is hard-coded for now.
func ID() uuid.UUID {
	return uuid.FromStringOrNil("743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb")
}
