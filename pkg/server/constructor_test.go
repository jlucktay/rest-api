package server_test

import (
	"reflect"
	"testing"

	"github.com/matryer/is"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/jlucktay/rest-api/pkg/storage/inmemory"
	"github.com/jlucktay/rest-api/pkg/storage/mongo"
)

func TestNew(t *testing.T) {
	testCases := map[string]struct {
		st       server.StorageType
		expected storage.PaymentStorage
	}{
		"In-memory storage (map); won't persist across app restarts": {
			st:       server.InMemory,
			expected: &inmemory.Storage{},
		},
		"Database storage (MongoDB); will persist across app restarts": {
			st:       server.Mongo,
			expected: &mongo.Storage{},
		},
	}
	for name, tC := range testCases {
		tC := tC // pin!
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			s := server.New(tC.st, true)
			is.Equal(reflect.TypeOf(s.Storage), reflect.TypeOf(tC.expected))
		})
	}
}
