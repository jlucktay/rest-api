package mongo

import (
	"log"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

func wrap(p storage.Payment, i ...mongoUUID) mongoWrapper {
	mw := mongoWrapper{
		Payment: p,
	}

	if len(i) > 0 {
		mw.UUID = i[0]
		return mw
	}

	newID, errID := uuid.NewV4()
	if errID != nil {
		log.Fatal(errID)
	}

	mw.UUID = mongoUUID(newID)

	return mw
}
