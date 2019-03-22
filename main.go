package main

import (
	"log"
	"net/http"
)

func main() {
	a := newAPIServer(Mongo)
	log.Fatal(http.ListenAndServe(":8080", a.router))
}
