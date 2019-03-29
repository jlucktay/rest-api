package inmemory

import (
	"sort"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	store map[uuid.UUID]storage.Payment
}

func (s *Storage) Init() error {
	s.store = make(map[uuid.UUID]storage.Payment)
	return nil
}

func (s *Storage) Create(p storage.Payment) (uuid.UUID, error) {
	newID := uuid.Must(uuid.NewV4())
	s.store[newID] = p
	return newID, nil
}

func (s *Storage) CreateSpecificID(id uuid.UUID, p storage.Payment) error {
	if _, exists := s.store[id]; exists {
		return &storage.AlreadyExistsError{ID: id}
	}
	s.store[id] = p
	return nil
}

func (s *Storage) Read(id uuid.UUID) (storage.Payment, error) {
	if p, exists := s.store[id]; exists {
		return p, nil
	}
	return storage.Payment{}, &storage.NotFoundError{ID: id}
}

func (s *Storage) ReadAll(rao storage.ReadAllOptions) (map[uuid.UUID]storage.Payment, error) {
	if rao.Limit == 0 {
		rao.Limit = storage.DefaultLimit
	}

	keys := make([]string, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)

	if uint(len(keys)) >= rao.Offset {
		keys = keys[rao.Offset:]
	} else {
		return map[uuid.UUID]storage.Payment{}, &storage.OffsetOutOfBoundsError{Offset: rao.Offset}
	}

	payments := make(map[uuid.UUID]storage.Payment)

	for i := uint(0); i < rao.Limit && i < uint(len(keys)); i++ {
		id := uuid.FromStringOrNil(keys[i])
		payments[id] = s.store[id]
	}

	return payments, nil
}

func (s *Storage) Update(id uuid.UUID, p storage.Payment) error {
	if _, exists := s.store[id]; exists {
		s.store[id] = p
		return nil
	}
	return &storage.NotFoundError{ID: id}
}

func (s *Storage) Delete(id uuid.UUID) error {
	if _, exists := s.store[id]; exists {
		delete(s.store, id)
		return nil
	}
	return &storage.NotFoundError{ID: id}
}
