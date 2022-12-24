package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chiefy/go-slack-utils/pkg/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

const (
	// SlackSigningSecret is the env var we look for that contains the signing secret for request validation
	SlackSigningSecret = "SLACK_SIGNING_SECRET"
)

var (
	signingSecret string
	debugMode     bool = false
)

func init() {
	signingSecret = os.Getenv(SlackSigningSecret)
	if signingSecret == "" {
		log.Fatalf("Could not find %s env var which is required", SlackSigningSecret)
	}
}

func main() {
	if os.Getenv("DEBUG") == "1" {
		debugMode = true
	}
	addr := "0.0.0.0:" + os.Getenv("PORT")

	r := mux.NewRouter()

	r.HandleFunc("/", MovieSearchHandler).Methods(http.MethodPost)
	r.HandleFunc("/lookup", MovieLookupHandler).Methods(http.MethodPost)

	if !debugMode {
		r.Use(middleware.ValidateTimestamp)
		r.Use(middleware.ValidateSlackRequest(signingSecret))
	}

	h := handlers.LoggingHandler(os.Stdout, r)

	srv := &http.Server{
		Handler:      h,
		Addr:         addr,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Printf("%s\nrunning on %s", GetVersion(), addr)
	log.Fatal(srv.ListenAndServe())
}
