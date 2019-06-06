// Package main launches an instance of the API server backed by Mongo for persistent storage.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jlucktay/rest-api/pkg/server"
)

func main() {
	var mongoHostname = flag.String("mongo-host", "localhost", "the hostname of the MongoDB server to connect to")
	flag.Parse()
	s := server.New(server.Mongo, *mongoHostname)
	log.Fatal(http.ListenAndServe(":8080", s.Router))
}
