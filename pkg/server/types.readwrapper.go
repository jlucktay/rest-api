package server

import (
	"net/http"

	"github.com/jlucktay/rest-api/internal/pkg/org"
	"github.com/jlucktay/rest-api/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

type ReadWrapper struct {
	Data  []paymentData    `json:"data"`
	Links ReadWrapperLinks `json:"links"`
}

type paymentData struct {
	Attributes     storage.Payment `json:"attributes"`
	ID             uuid.UUID       `json:"id"`
	OrganisationID uuid.UUID       `json:"organisation_id"`
	Type           string          `json:"type"`
	Version        int             `json:"version"`
}

type ReadWrapperLinks struct {
	Self string `json:"self"`
}

func (rw *ReadWrapper) Init(r *http.Request) {
	rw.Data = make([]paymentData, 0)
	rw.Links = ReadWrapperLinks{
		Self: r.URL.String(),
	}
}

func (rw *ReadWrapper) addPayment(id uuid.UUID, p storage.Payment) {
	newPD := &paymentData{
		Attributes:     p,
		ID:             id,
		OrganisationID: org.ID(),
		Type:           "Payment",
		Version:        0,
	}

	rw.Data = append(rw.Data, *newPD)
}
