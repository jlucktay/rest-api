package mongo

import (
	"log"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

func wrap(p storage.Payment, i ...uuid.UUID) mongoWrapper {
	mw := mongoWrapper{
		Payment: p,
		UUID:    &mongoUUID{},
	}

	if len(i) > 0 {
		mw.UUID.UUID = i[0]
		return mw
	}

	newID, errID := uuid.NewV4()
	if errID != nil {
		log.Fatal(errID)
	}

	mw.UUID.UUID = newID

	return mw
}
