package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type dummyStorage struct {
	store map[uuid.UUID]Payment
}

func (ds *dummyStorage) Init() error {
	ds.store = make(map[uuid.UUID]Payment)
	return nil
}

func (ds *dummyStorage) Create(p Payment) (uuid.UUID, error) {
	newId, errNew := uuid.NewV4()
	if errNew != nil {
		return uuid.Nil, errNew
	}
	ds.store[newId] = p
	return newId, nil
}

func (ds *dummyStorage) Read(id uuid.UUID) (Payment, error) {
	if p, ok := ds.store[id]; ok {
		return p, nil
	}
	return Payment{}, &ReadError{id}
}

func (ds *dummyStorage) Update(id uuid.UUID, p Payment) error {
	if _, ok := ds.store[id]; ok {
		ds.store[id] = p
		return nil
	}
	return fmt.Errorf("Payment ID '%s' not found.", id)
}

func (ds *dummyStorage) Delete(id uuid.UUID) error {
	if _, ok := ds.store[id]; ok {
		delete(ds.store, id)
		return nil
	}
	return fmt.Errorf("Payment ID '%s' not found.", id)
}