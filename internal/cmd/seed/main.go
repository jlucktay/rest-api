package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/jlucktay/rest-api/pkg/server"
	"github.com/jlucktay/rest-api/pkg/storage/mongo"
)

func main() {
	// Empty out current MongoDB collection.
	m := mongo.New()
	if errInit := m.Initialise(); errInit != nil {
		log.Fatal("error initialising MongoDB the first time:", errInit)
	}
	if errTerm := m.Terminate(true); errTerm != nil {
		log.Fatal("error emptying MongoDB:", errTerm)
	}

	// Renew our MongoDB client's connection to the database.
	m = mongo.New()
	if errInitAgain := m.Initialise(); errInitAgain != nil {
		log.Fatal("error initialising MongoDB the second time:", errInitAgain)
	}

	// Get the list of payments from Mockbin.
	seedUrl, errParse := url.Parse("http://mockbin.org/bin/41ca3269-d8c4-4063-9fd5-f306814ff03f")
	if errParse != nil {
		log.Fatal("error parsing URL:", errParse)
	}
	resp, errGet := getResponse(*seedUrl)
	if errGet != nil {
		log.Fatal("error getting response:", errGet)
	}

	respBytes, errRead := ioutil.ReadAll(resp)
	if errRead != nil {
		log.Fatal("error reading bytes:", errRead)
	}

	// Parse the payments out into a struct.
	rw := &server.ReadWrapper{}

	if errUm := json.Unmarshal(respBytes, rw); errUm != nil {
		log.Fatal("error unmarshaling response:", errUm)
	}

	// Create the individual payments in storage, and preserve their UUIDs.
	for _, pd := range rw.Data {
		errCreate := m.CreateSpecificID(pd.ID, pd.Attributes)
		if errCreate != nil {
			log.Fatal("error creating payment:", errCreate)
		}
		fmt.Printf("Added payment with ID '%v'.\n", pd.ID)
	}
}
