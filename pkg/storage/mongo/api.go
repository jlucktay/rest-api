package mongo

import (
	"errors"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

func (s *Storage) Init() error {
	return errors.New("not yet implemented")
}

func (s *Storage) Create(p storage.Payment) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not yet implemented")
}

func (s *Storage) CreateSpecificID(id uuid.UUID, p storage.Payment) error {
	return errors.New("not yet implemented")
}

func (s *Storage) Read(id uuid.UUID) (storage.Payment, error) {
	return storage.Payment{}, errors.New("not yet implemented")
}

func (s *Storage) ReadAll(rao storage.ReadAllOptions) (map[uuid.UUID]storage.Payment, error) {
	return make(map[uuid.UUID]storage.Payment), errors.New("not yet implemented")
}

func (s *Storage) Update(id uuid.UUID, p storage.Payment) error {
	return errors.New("not yet implemented")
}

func (s *Storage) Delete(id uuid.UUID) error {
	return errors.New("not yet implemented")
}
