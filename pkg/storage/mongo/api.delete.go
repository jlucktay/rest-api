package mongo

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

func (s *Storage) Delete(id uuid.UUID) error {
	return errors.New("not yet implemented")
}
