// Package main launches an instance of the API server backed by Mongo for persistent storage.
package main

import (
	"flag"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"go.jlucktay.dev/rest-api/pkg/server"
)

const defaultMongo = "localhost"

func main() {
	mongoHostname := flag.String("mongo-host", defaultMongo, "the hostname of the MongoDB server to connect to")

	flag.Parse()

	// If the flag has not been given, look for an environment variable.
	if envMongo, found := os.LookupEnv("MONGO_HOST"); found && *mongoHostname == defaultMongo {
		*mongoHostname = envMongo
	}

	s := server.New(server.Mongo, false, *mongoHostname)
	log.Fatal(http.ListenAndServe(":8080", s.Router))
}
