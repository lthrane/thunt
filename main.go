package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"

	"github.com/gorilla/mux"
)

var client *datastore.Client

func main() {
	ctx := context.Background()

	projectID := "thunt-269016"
	var err error
	client, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create Datastore client: %v", err)
	}

	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout)

	router := mux.NewRouter().StrictSlash(true)
	register(router)

	Info.Printf("starting server")
	log.Fatal(http.ListenAndServe(":8080", router))
}
