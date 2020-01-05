package main

import (
	"github.com/chiefy/go-slack-utils/pkg/middleware"
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
	signingSecret string
)

func init() {
	signingSecret = os.Getenv(SlackSigningSecret)
	if signingSecret == "" {
		log.Fatalf("Could not find %s env var which is required", SlackSigningSecret)
	}
}

func main() {
	addr := "127.0.0.1:" + os.Getenv("PORT")

	r := mux.NewRouter()
	r.HandleFunc("/", MovieSearchHandler).Methods(http.MethodPost)
	r.Use(middleware.ValidateTimestamp)
	r.Use(middleware.ValidateSlackRequest(signingSecret))

	r.HandleFunc("/lookup", MovieLookupHandler).Methods(http.MethodPost)
	r.Use(middleware.ValidateTimestamp)
	r.Use(middleware.ValidateSlackRequest(signingSecret))

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(GetVersion())
	log.Fatal(srv.ListenAndServe())
}
