package server

import (
	"github.com/go-chi/chi"
)

func (s *Server) setupRoutes() {
	// s.Router.Use() // TODO middleware?

	// RESTy routes for "payments" resource
	s.Router.Route("/payments", func(r chi.Router) {
		r.Post("/", s.createPayments())
		r.Get("/", s.readPayments())

		r.Route("/{id}", func(r chi.Router) {
			r.Post("/", s.createPaymentByID())
			r.Get("/", s.readPaymentByID())
			r.Put("/", s.updatePaymentByID())
			r.Delete("/", s.deletePaymentByID())
		})
	})
}

/*
// TODO - paginate?
r.Route("/articles", func(r chi.Router) {
	r.With(paginate).Get("/", ListArticles)
})
*/
