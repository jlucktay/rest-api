package server

import (
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func (s *Server) setupRoutes() {
	log.Debug("Setting up routes...")

	// RESTy routes for v1 "payments" resource
	s.Router.Route("/v1", func(r chi.Router) {
		r.Route("/payments", func(r chi.Router) {
			r.Post("/", s.createPayments())
			r.Get("/", s.readPayments())

			r.Route("/{id}", func(r chi.Router) {
				r.Post("/", s.createPaymentByID())
				r.Get("/", s.readPaymentByID())
				r.Put("/", s.updatePaymentByID())
				r.Delete("/", s.deletePaymentByID())
			})
		})
	})

	log.Debug("Set up routes.")
}
