package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/chiefy/go-slack-utils/pkg/middleware"
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

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)

		slog.Info("request received", userAttrs, requestAttrs)
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	if os.Getenv("DEBUG") == "1" {
		log.Println("Running in DEBUG mode!")
		debugMode = true
	}

	addr := "0.0.0.0:" + getEnv("PORT", "5555")

	mux := http.NewServeMux()

	mux.HandleFunc("/", MovieSearchHandler)
	mux.HandleFunc("/lookup", MovieLookupHandler)

	mw := logRequest(mux)

	if !debugMode {
		mw = middleware.ValidateTimestamp(mw)
		mw = middleware.ValidateSlackRequest(signingSecret)(mw)
	}

	srv := &http.Server{
		Handler:      mw,
		Addr:         addr,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("%s\nrunning on %s", GetVersion(), addr)
	log.Fatal(srv.ListenAndServe())
}
