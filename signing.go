package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	slackSigHeader          = "X-Slack-Signature"
	slackSigHeaderTimestamp = "X-Slack-Request-Timestamp"
)

func checkTimestampMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !validTS(getTS(r)) {
			log.Println("invalid timestamp")
			http.Error(w, "Invalid Timestamp", http.StatusBadRequest)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func checkSlackSigningMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("bad request body %s", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		incomingSig := []byte(fmt.Sprintf("v0:%d:%s", getTS(r), string(bodyData)))
		slackSig, _ := hex.DecodeString(strings.TrimPrefix(r.Header.Get(slackSigHeader), "v0="))

		if !validMAC(incomingSig, slackSig, signingSecret) {
			log.Println("HMAC error - did not match")
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyData))
			next.ServeHTTP(w, r)
		}
	})
}

func getTS(r *http.Request) int64 {
	ts, err := strconv.ParseInt(r.Header.Get(slackSigHeaderTimestamp), 10, 64)
	if err != nil {
		log.Println("could not get timestamp", err)
	}
	return ts
}

func validTS(ts int64) bool {
	m, _ := time.ParseDuration("2m")
	return time.Since(time.Unix(ts, 0)) < m
}

// validMAC reports whether messageMAC is a valid HMAC tag for message.
func validMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
