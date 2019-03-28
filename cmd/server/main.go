package main

import (
	"log"
	"net/http"

	_ "github.com/jlucktay/rest-api"
)

func main() {
	a := newAPIServer(Mongo)
	log.Fatal(http.ListenAndServe(":8080", a.router))
}
