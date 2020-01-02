package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Message is the Slack message event.
// see https://api.slack.com/docs/message-formatting
type Message struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// Log is the user log for the barometer. This also serves as
// the schema for the BigQuery table
type Log struct {
}

// BurnoutBarometer takes a log message from a Slack slash command and stores
// it into BigQuery as streaming insert.
func BurnoutBarometer(w http.ResponseWriter, r *http.Request) {
	setup(r.Context())
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", 405)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Couldn't parse form", 400)
		log.Fatalf("ParseForm: %v", err)
	}

	if err := verifyWebHook(r.Form); err != nil {
		log.Fatalf("verifyWebHook: %v", err)
	}

	if len(r.Form["text"]) == 0 {
		log.Fatalf("Empty text in form")
	}

	bbResponse, err := storeMessage(r.Form["text"][0])
	if err != nil {
		log.Fatalf("storeMessage: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(bbResponse); err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
}

func verifyWebHook(form url.Values) error {
	t := form.Get("token")
	if len(t) == 0 {
		return fmt.Errorf("empty form token")
	}

	if t != config.Token {
		return fmt.Errorf("invalid request/credentials: %q", t[0])
	}

	return nil
}

func storeMessage(msg string) (*Message, error) {
	// TODO: Store message in BigQuery as a streaming insert
}
