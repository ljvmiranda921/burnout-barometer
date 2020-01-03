package barometer

import (
	"fmt"
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
	}, "", nil
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
