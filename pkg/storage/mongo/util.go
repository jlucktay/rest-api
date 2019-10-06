package mongo

import (
	"github.com/gofrs/uuid"

	"github.com/jlucktay/rest-api/pkg/storage"
)

func wrap(p storage.Payment, i ...mongoUUID) mongoWrapper {
	mw := mongoWrapper{Payment: p}

	if len(i) > 0 {
		mw.UUID = i[0]
		return mw
	}

	newID := uuid.Must(uuid.NewV4())
	mw.UUID = mongoUUID(newID)

	return mw
}
