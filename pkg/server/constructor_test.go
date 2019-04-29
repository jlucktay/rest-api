package server_test

import (
	"reflect"
	"testing"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage"
	"github.com/jlucktay/rest-api/pkg/storage/inmemory"
	"github.com/jlucktay/rest-api/pkg/storage/mongo"
	"github.com/matryer/is"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		desc     string
		st       server.StorageType
		expected storage.PaymentStorage
	}{
		{
			desc:     "In-memory storage (map); won't persist across app restarts",
			st:       server.InMemory,
			expected: &inmemory.Storage{},
		},
		{
			desc:     "Database storage (MongoDB); will persist across app restarts",
			st:       server.Mongo,
			expected: &mongo.Storage{},
		},
	}
	for _, tC := range testCases {
		tC := tC // pin!
		t.Run(tC.desc, func(t *testing.T) {
			i := is.New(t)
			s := server.New(tC.st, true)
			i.Equal(reflect.TypeOf(s.Storage), reflect.TypeOf(tC.expected))
		})
	}
}
