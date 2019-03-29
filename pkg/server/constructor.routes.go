package server

func (s *Server) setupRoutes() {
	s.Router.HandleMethodNotAllowed = true

	// C
	s.Router.POST("/payments", s.createPayments())
	s.Router.POST("/payments/:id", s.createPaymentByID())

	// R
	s.Router.GET("/payments", s.readPayments())
	s.Router.GET("/payments/:id", s.readPaymentByID())

	// U
	s.Router.PUT("/payments/:id", s.updatePaymentByID())

	// D
	s.Router.DELETE("/payments/:id", s.deletePaymentByID())
}
