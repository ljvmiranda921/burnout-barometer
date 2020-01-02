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

	"cloud.google.com/go/bigquery"
)

// Log is the user log for the barometer. This also serves as
// the schema for the BigQuery table.
type Log struct {
	Timestamp  time.Time
	UserID     string
	LogMeasure int
	Notes      string
}

// Save implements the ValueSaver interface.
func (i *Log) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"timestamp":   i.Timestamp,
		"user_id":     i.UserID,
		"log_measure": i.LogMeasure,
		"notes":       i.Notes,
	}
}

// FormatReply prepares the Slack message as a response to a slash command.
func (i *Log) FormatReply() (*Message, error) {
	attach := Attachment{
		Color: "#ef4631",
		Title: "Burnout Barometer",
		Text:  fmt.Sprintf("Acknowledged"),
	}

	message := &Message{
		ResponseType: "ephemeral",
		Text:         fmt.Sprintf("Received: %d (%s)", i.LogMeasure, i.Notes),
		Attachments:  []Attachment{attach},
	}

	return message, nil
}

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

	ts, err := getTimestamp(date, config.Area)
	if err != nil {
		log.Fatalf("error in getTimestamp: %v", err)
	}

	measure, err := strconv.Atoi(l[0])
	if err != nil {
		log.Fatalf("error in strconv: %v", err)
	}

	logMsg := &Log{
		Timestamp:  ts,
		UserID:     userID,
		LogMeasure: measure,
		Notes:      l[1],
	}

	// TODO: Store message in BigQuery as a streaming insert

	return logMsg.FormatReply()
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
		return time.Time{}, err
	}
	loc, err := time.LoadLocation(config.Area)
	if err != nil {
		log.Fatalf("cannot find location: %s", area)
		return time.Time{}, err
	}

	return time.Unix(i, 0).In(loc), nil
}
