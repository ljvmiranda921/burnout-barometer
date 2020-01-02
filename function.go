package barometer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// BurnoutBarometer takes a log message from a Slack slash command and stores
// it into BigQuery as a streaming insert.
func BurnoutBarometer(w http.ResponseWriter, r *http.Request) {

	log.Printf("request received")

	// Setup application variables
	if err := setup(r.Context()); err != nil {
		log.Fatalf("setup: %v", err)
	}

	// Validate request and parse the submitted form
	if r.Method != "POST" {
		http.Error(w, "only POST requests are accepted", 405)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Couldn't parse form", 400)
		log.Fatalf("error in ParseForm: %v", err)
	}

	if err := verifyWebHook(r.Form); err != nil {
		log.Fatalf("error in verifyWebHook: %v", err)
	}

	if len(r.Form["text"]) == 0 {
		log.Fatalf("empty text in form")
	}

	// Store the message and timestamp to BigQuery
	req := &Request{
		Text:      r.Form["text"][0],
		UserID:    r.Form["user_id"][0],
		Timestamp: r.Header.Get("X-Slack-Request-Timestamp"),
		Area:      config.Area,
	}
	resp, err := req.Process()
	if err != nil {
		log.Fatalf("error in storeMessage: %v", err)
	}

	// Send reply back to Slack
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Fatalf("error in json.Marshal: %v", err)
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

// Request defines the common form parameters when a slash command is invoked.
type Request struct {
	Text      string
	UserID    string
	Timestamp string
	Area      string
	LogMsg    Log
}

// Process parses the request and stores to BigQuery.
func (r *Request) Process() (*Message, error) {
	m, notes, err := r.ParseMessage()
	if err != nil {
		log.Fatalf("error in parseLogMessage: %v", err)
		return nil, err
	}

	ts, err := r.GetTimestamp()
	if err != nil {
		log.Fatalf("error in getTimestamp: %v", err)
	}

	measure, err := strconv.Atoi(m)
	if err != nil {
		log.Fatalf("error in strconv: %v", err)
	}

	r.LogMsg = Log{
		Timestamp:  ts,
		UserID:     r.UserID,
		LogMeasure: measure,
		Notes:      notes,
	}

	// TODO: Store message in BigQuery as a streaming insert

	return r.LogMsg.FormatReply()
}

// ParseMessage extracts the barometer measure and notes from the form text.
func (r *Request) ParseMessage() (string, string, error) {
	list := strings.Fields(r.Text)
	measure := list[0]
	notes := strings.Join(list[1:], " ")
	return measure, notes, nil
}

// GetTimestamp obtains the timestamp value from the request.
func (r *Request) GetTimestamp() (time.Time, error) {
	i, err := strconv.ParseInt(r.Timestamp, 10, 64)
	if err != nil {
		log.Fatalf("cannot parse timestamp %s: %v", r.Timestamp, err)
		return time.Time{}, err
	}
	loc, err := time.LoadLocation(r.Area)
	if err != nil {
		log.Fatalf("cannot find location: %s", r.Area)
		return time.Time{}, err
	}

	return time.Unix(i, 0).In(loc), nil
}
