package inmemory

import (
	"sort"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

type InMemoryStorage struct {
	store map[uuid.UUID]storage.Payment
}

const defaultLimit = 10

func (ims *InMemoryStorage) Init() error {
	ims.store = make(map[uuid.UUID]storage.Payment)
	return nil
}

func (ims *InMemoryStorage) Create(p storage.Payment) (uuid.UUID, error) {
	newID := uuid.Must(uuid.NewV4())
	ims.store[newID] = p
	return newID, nil
}

func (ims *InMemoryStorage) createSpecificID(id uuid.UUID, p storage.Payment) error {
	if _, exists := ims.store[id]; exists {
		return &AlreadyExistsError{id}
	}
	ims.store[id] = p
	return nil
}

func (ims *InMemoryStorage) Read(id uuid.UUID) (storage.Payment, error) {
	if p, exists := ims.store[id]; exists {
		return p, nil
	}
	return storage.Payment{}, &NotFoundError{id}
}

func (ims *InMemoryStorage) ReadAll(rao storage.ReadAllOptions) (map[uuid.UUID]storage.Payment, error) {
	if rao.Limit == 0 {
		rao.Limit = defaultLimit
	}

	keys := make([]string, 0, len(ims.store))
	for k := range ims.store {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)

	if uint(len(keys)) >= rao.Offset {
		keys = keys[rao.Offset:]
	} else {
		return map[uuid.UUID]storage.Payment{}, &OffsetOutOfBoundsError{rao.Offset}
	}

	payments := make(map[uuid.UUID]storage.Payment)

	for i := uint(0); i < rao.Limit && i < uint(len(keys)); i++ {
		id := uuid.FromStringOrNil(keys[i])
		payments[id] = ims.store[id]
	}

	return payments, nil
}

func (ims *InMemoryStorage) Update(id uuid.UUID, p storage.Payment) error {
	if _, exists := ims.store[id]; exists {
		ims.store[id] = p
		return nil
	}
	return &NotFoundError{id}
}

func (ims *InMemoryStorage) Delete(id uuid.UUID) error {
	if _, exists := ims.store[id]; exists {
		delete(ims.store, id)
		return nil
	}
	return &NotFoundError{id}
}
