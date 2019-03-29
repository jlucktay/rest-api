// Package org represents an organisation that makes payment transactions.
package org

import (
	uuid "github.com/satori/go.uuid"
)

// ID is hard-coded for now.
func ID() uuid.UUID {
	return uuid.FromStringOrNil("a6781162-0f4f-429c-aca1-ac7a0cff4edf")
}
