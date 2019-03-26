package main

import (
	"net/http"

	"github.com/jlucktay/rest-api/internal/pkg/org"
	uuid "github.com/satori/go.uuid"
)

type readWrapper struct {
	Data  []paymentData    `json:"data"`
	Links readWrapperLinks `json:"links"`
}

type paymentData struct {
	Attributes     Payment   `json:"attributes"`
	ID             uuid.UUID `json:"id"`
	OrganisationID uuid.UUID `json:"organisation_id"`
	Type           string    `json:"type"`
	Version        int       `json:"version"`
}

type readWrapperLinks struct {
	Self string `json:"self"`
}

func (rw *readWrapper) init(r *http.Request) {
	rw.Data = make([]paymentData, 0)
	rw.Links = readWrapperLinks{
		Self: r.URL.String(),
	}
}

func (rw *readWrapper) addPayment(id uuid.UUID, p Payment) {
	newPD := &paymentData{
		Attributes:     p,
		ID:             id,
		OrganisationID: org.ID(),
		Type:           "Payment",
		Version:        0,
	}

	rw.Data = append(rw.Data, *newPD)
}