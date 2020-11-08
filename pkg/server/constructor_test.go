package server_test

import (
	"reflect"
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/rest-api/pkg/server"
	"go.jlucktay.dev/rest-api/pkg/storage"
	"go.jlucktay.dev/rest-api/pkg/storage/inmemory"
	jramongo "go.jlucktay.dev/rest-api/pkg/storage/mongo"
)

func TestNew(t *testing.T) {
	testCases := map[string]struct {
		st         server.StorageType
		expected   storage.PaymentStorage
		connString string
	}{
		"In-memory storage (map); won't persist across app restarts": {
			st:       server.InMemory,
			expected: &inmemory.Storage{},
		},
		"Database storage (MongoDB); will persist across app restarts": {
			st:         server.Mongo,
			expected:   &jramongo.Storage{},
			connString: directConnString,
		},
	}

	for name, tC := range testCases {
		tC := tC // pin!

		t.Run(name, func(_ *testing.T) {
			var s *server.Server
			if tC.connString != "" {
				s = server.New(tC.st, true, tC.connString)
			} else {
				s = server.New(tC.st, true)
			}
			is.Equal(reflect.TypeOf(s.Storage), reflect.TypeOf(tC.expected)) // wrong storage type
		})
	}
}
