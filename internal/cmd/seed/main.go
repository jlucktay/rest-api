package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage/mongo"
	"github.com/sirupsen/logrus"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Exiting!")
		os.Exit(0)
	}()

	fmt.Println("Continuing will delete ALL payment records in MongoDB (database: rest-api, collection: payments")
	fmt.Print("Press 'Enter' to continue, or CTRL+C to cancel...")
	_, errRead := bufio.NewReader(os.Stdin).ReadBytes('\n')
	theShowMustGoOn("error reading from stdin", errRead)

	// Empty out current MongoDB collection.
	m := mongo.New()
	errInit := m.Initialise()
	theShowMustGoOn("error initialising MongoDB the first time", errInit)

	errTerm := m.Terminate(true)
	theShowMustGoOn("error emptying MongoDB", errTerm)

	// Renew our MongoDB client's connection to the database.
	m = mongo.New()
	errInitAgain := m.Initialise()
	theShowMustGoOn("error initialising MongoDB the second time", errInitAgain)

	// Get the list of payments from Mockbin.
	seedURL, errParse := url.Parse("http://mockbin.org/bin/41ca3269-d8c4-4063-9fd5-f306814ff03f")
	theShowMustGoOn("error parsing URL", errParse)

	resp, errGet := getResponse(*seedURL)
	theShowMustGoOn("error getting response", errGet)

	respBytes, errRead := ioutil.ReadAll(resp)
	theShowMustGoOn("error reading bytes", errRead)

	// Parse the payments out into a struct.
	rw := &server.ReadWrapper{}
	errUm := json.Unmarshal(respBytes, rw)
	theShowMustGoOn("error unmarshaling response", errUm)

	// Create the individual payments in storage, and preserve their UUIDs.
	for _, pd := range rw.Data {
		errCreate := m.CreateSpecificID(pd.ID, pd.Attributes)
		theShowMustGoOn("error creating payment", errCreate)
		fmt.Printf("Added payment with ID '%v'.\n", pd.ID)
	}
}

func theShowMustGoOn(s string, e error) {
	if e != nil {
		logrus.Fatalf("%s: %s\n", s, e)
	}
}
