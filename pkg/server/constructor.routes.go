package server

func (s *Server) setupRoutes() {
	// C
	s.Router.Post("/payments", s.createPayments())
	s.Router.Post("/payments/{id}", s.createPaymentByID())

	// R
	s.Router.Get("/payments", s.readPayments())
	s.Router.Get("/payments/{id}", s.readPaymentByID())

	// U
	s.Router.Put("/payments/{id}", s.updatePaymentByID())

	// D
	s.Router.Delete("/payments/{id}", s.deletePaymentByID())
}
