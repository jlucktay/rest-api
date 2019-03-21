package main

import (
	"sort"

	uuid "github.com/satori/go.uuid"
)

type inMemoryStorage struct {
	store map[uuid.UUID]Payment
}

const defaultLimit = 100

func (ims *inMemoryStorage) Init() error {
	ims.store = make(map[uuid.UUID]Payment)
	return nil
}

func (ims *inMemoryStorage) Create(p Payment) (uuid.UUID, error) {
	newId := uuid.Must(uuid.NewV4())
	ims.store[newId] = p
	return newId, nil
}

func (ims *inMemoryStorage) createSpecificId(id uuid.UUID, p Payment) error {
	if _, exists := ims.store[id]; exists {
		return &AlreadyExistsError{id}
	}
	ims.store[id] = p
	return nil
}

func (ims *inMemoryStorage) Read(id uuid.UUID) (Payment, error) {
	if p, exists := ims.store[id]; exists {
		return p, nil
	}
	return Payment{}, &NotFoundError{id}
}

func (ims *inMemoryStorage) ReadAll(rao ReadAllOptions) ([]Payment, error) {
	if rao.limit == 0 {
		rao.limit = defaultLimit
	}

	var keys []string
	for k := range ims.store {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)

	if uint(len(keys)) >= rao.offset {
		keys = keys[rao.offset:]
	} else {
		return []Payment{}, &OffsetOutOfBounds{rao.offset}
	}

	payments := make([]Payment, 0, rao.limit)
	for i := uint(0); i < rao.limit && i < uint(len(keys)); i++ {
		payments = append(payments, ims.store[uuid.FromStringOrNil(keys[i])])
	}

	return payments, nil
}

func (ims *inMemoryStorage) Update(id uuid.UUID, p Payment) error {
	if _, exists := ims.store[id]; exists {
		ims.store[id] = p
		return nil
	}
	return &NotFoundError{id}
}

func (ims *inMemoryStorage) Delete(id uuid.UUID) error {
	if _, exists := ims.store[id]; exists {
		delete(ims.store, id)
		return nil
	}
	return &NotFoundError{id}
}
