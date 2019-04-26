// Package main launches an instance of the API server backed by Mongo for persistent storage.
package main

import (
	"flag"
	"net/http"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	var logDebug = flag.Bool("debug", false, "enable debug logging")
	var mongoHostname = flag.String("mongo-host", "localhost", "the hostname of the MongoDB server to connect to")
	flag.Parse()

	if *logDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})

	s := server.New(server.Mongo, *mongoHostname)
	logrus.Fatal(http.ListenAndServe(":8080", s.Router))
}
