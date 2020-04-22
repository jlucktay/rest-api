package inmemory

import (
	"sort"

	"github.com/gofrs/uuid"

	"go.jlucktay.dev/rest-api/pkg/storage"
)

// Storage is an in-memory storage system that uses a map to store Payment structs indexed by UUID.
type Storage struct {
	store map[uuid.UUID]storage.Payment
}

// Initialise will initialise the internal map.
func (s *Storage) Initialise() error {
	s.store = make(map[uuid.UUID]storage.Payment)
	return nil
}

// Terminate will terminate the internal map by setting the internal store to nil, so that the garbage collector will
// pick up the map's contents.
// It does not bother checking the true/false argument for whether or not to destroy the stored data, due to the very
// nature of its implementation.
func (s *Storage) Terminate(...bool) error {
	s.store = nil
	return nil
}

// Create will add the given Payment to the store, and return its UUID.
func (s *Storage) Create(p storage.Payment) (uuid.UUID, error) {
	newID := uuid.Must(uuid.NewV4())
	s.store[newID] = p

	return newID, nil
}

// CreateSpecificID will create the given Payment with that specific UUID, rather than generating its own.
// Intended for use with testing, and not in production.
func (s *Storage) CreateSpecificID(id uuid.UUID, p storage.Payment) error {
	if _, exists := s.store[id]; exists {
		return &storage.AlreadyExistsError{ID: id}
	}

	s.store[id] = p

	return nil
}

// Read will attempt to find the Payment linked to the given UUID.
func (s *Storage) Read(id uuid.UUID) (storage.Payment, error) {
	if p, exists := s.store[id]; exists {
		return p, nil
	}

	return storage.Payment{}, &storage.NotFoundError{ID: id}
}

// ReadAll will read all available Payments, up to the limit which is either 1) explicit in the given ReadAllOptions
// struct, or 2) the default limit which is also exported by this package.
// The Payments will be returned in order based on their respective UUIDs.
// The entire collection of Payments can be paginated through via the Offset property of the ReadAllOptions struct, and
// by default when such an Offset is not explicitly specified, the returned collection will start from the beginning of
// the store.
func (s *Storage) ReadAll(rao storage.ReadAllOptions) (map[uuid.UUID]storage.Payment, error) {
	// Set limit from options or default constant.
	if rao.Limit == 0 {
		rao.Limit = storage.DefaultLimit
	}

	// Get all keys and sort in order.
	keys := make([]string, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k.String())
	}

	sort.Strings(keys)

	// Truncate keys slice if applicable.
	if uint(len(keys)) >= rao.Offset {
		keys = keys[rao.Offset:]
	} else {
		return map[uuid.UUID]storage.Payment{}, &storage.OffsetOutOfBoundsError{Offset: rao.Offset}
	}

	// Build result set.
	payments := make(map[uuid.UUID]storage.Payment)

	for i := uint(0); i < rao.Limit && i < uint(len(keys)); i++ {
		id := uuid.FromStringOrNil(keys[i])
		payments[id] = s.store[id]
	}

	return payments, nil
}

// Update will overwrite the Payment currently stored at the given UUID with the given Payment.
func (s *Storage) Update(id uuid.UUID, p storage.Payment) error {
	if _, exists := s.store[id]; exists {
		s.store[id] = p
		return nil
	}

	return &storage.NotFoundError{ID: id}
}

// Delete will remove the Payment at the given UUID from the internal store.
func (s *Storage) Delete(id uuid.UUID) error {
	if _, exists := s.store[id]; exists {
		delete(s.store, id)
		return nil
	}

	return &storage.NotFoundError{ID: id}
}
