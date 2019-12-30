package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// SlackSigningSecret is the env var we look for that contains the signing secret for request validation
	SlackSigningSecret = "SLACK_SIGNING_SECRET"
)

var (
	signingSecret []byte
)

func init() {
	ss := os.Getenv(SlackSigningSecret)
	if ss == "" {
		log.Fatalf("Could not find %s env var which is required", SlackSigningSecret)
	}
	signingSecret = []byte(ss)
}

func main() {
	addr := "127.0.0.1:" + os.Getenv("PORT")

	r := mux.NewRouter()
	r.HandleFunc("/", MovieHandler).Methods(http.MethodPost)
	r.Use(checkTimestampMiddleware)
	r.Use(checkSlackSigningMiddleware)

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
