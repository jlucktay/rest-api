package mongo

import (
	"log"

	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

func wrap(p storage.Payment, i ...uuid.UUID) mongoWrapper {
	mw := mongoWrapper{
		Payment: p,
		UUID:    "",
	}

	if len(i) > 0 {
		mw.UUID = i[0].String()
		return mw
	}

	newID, errID := uuid.NewV4()
	if errID != nil {
		log.Fatal(errID)
	}

	mw.UUID = newID.String()

	return mw
}
