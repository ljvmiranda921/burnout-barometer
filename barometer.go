package function

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

// Message is the Slack message event.
// see https://api.slack.com/docs/message-formatting
type Message struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

// Attachment defines the message output after running the slash command
type Attachment struct {
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
}

// Log is the user log for the barometer. This also serves as
// the schema for the BigQuery table
type Log struct {
	Timestamp  time.Time
	UserID     string
	LogMeasure int
	Notes      string
}

// BurnoutBarometer takes a log message from a Slack slash command and stores
// it into BigQuery as streaming insert.
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
	resp, err := storeMessage(
		r.Form["text"][0],
		r.Form["user_id"][0],
		r.Header.Get("X-Slack-Request-Timestamp"),
	)
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

func storeMessage(text, userID, date string) (*Message, error) {
	// Prepare message to be sent to BigQuery
	l, err := parseLogMessage(text)
	if err != nil {
		log.Fatalf("error in parseLogMessage: %v", err)
		return nil, err
	}

	logMsg := &Log{
		Timestamp:  getTimestamp(date, config.Area),
		UserID:     userID,
		LogMeasure: strconv.Atoi(l[0]),
		Notes:      l[1],
	}

	// TODO: Store message in BigQuery as a streaming insert
	log.Printf(logMsg)

	return FormatReply(logMsg), nil
}

func parseLogMessage(m string) ([]string, error) {
	list := strings.Fields(m)
	measure := list[0]
	notes := strings.Join(list[1:], " ")
	return []string{measure, notes}, nil
}

func getTimestamp(t string, area string) (time.Time, error) {
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		log.Fatalf("cannot parse timestamp %s: %v", t, err)
		return nil, err
	}
	loc, err := time.LoadLocation(config.Area)
	if err != nil {
		log.Fatalf("cannot find location: %s", area)
		return nil, err
	}

	return time.Unix(i, 0).In(loc), nil
}
