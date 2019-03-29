package main

import (
	"log"
	"net/http"

	"github.com/jlucktay/rest-api/pkg/server"
)

func main() {
	s := server.New(server.Mongo)
	log.Fatal(http.ListenAndServe(":8080", s.Router))
}
