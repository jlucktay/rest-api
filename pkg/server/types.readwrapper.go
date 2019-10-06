package server

import (
	"github.com/gofrs/uuid"

	"github.com/jlucktay/rest-api/pkg/org"
	"github.com/jlucktay/rest-api/pkg/storage"
)

// ReadWrapper adds some extra information around Payment structs that are read from the API.
type ReadWrapper struct {
	Data  []PaymentData    `json:"data"`
	Links ReadWrapperLinks `json:"links"`
}

type PaymentData struct {
	Attributes     storage.Payment `json:"attributes"`
	ID             uuid.UUID       `json:"id"`
	OrganisationID uuid.UUID       `json:"organisation_id"`
	Type           string          `json:"type"`
	Version        int             `json:"version"`
}

type ReadWrapperLinks struct {
	Self string `json:"self"`
}

// NewWrapper will return a new ReadWrapper.
func NewWrapper(s string) *ReadWrapper {
	rw := new(ReadWrapper)
	rw.Data = make([]PaymentData, 0)
	rw.Links = ReadWrapperLinks{
		Self: s,
	}
	return rw
}

// AddPayment will add a Payment with some other boilerplate attributes to a ReadWrapper.
func (rw *ReadWrapper) AddPayment(id uuid.UUID, p storage.Payment) {
	newPD := &PaymentData{
		Attributes:     p,
		ID:             id,
		OrganisationID: org.ID(),
		Type:           "Payment",
		Version:        0,
	}

	rw.Data = append(rw.Data, *newPD)
}
