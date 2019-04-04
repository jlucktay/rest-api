package mongo

import (
	"errors"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

func (s *Storage) Update(id uuid.UUID, p storage.Payment) error {
	return errors.New("not yet implemented")
}
