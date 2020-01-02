package barometer

import "fmt"

// FormatReply prepares the Slack message as a response to a slash command
func FormatReply(m *Log) (*Message, error) {

	attach := Attachment{
		Color: "#ef4631",
		Title: "Burnout Barometer",
		Text:  fmt.Sprintf("Acknowledged"),
	}

	message := &Message{
		ResponseType: "ephemeral",
		Text:         fmt.Sprintf("Received: %s (%s)", m.LogMeasure, m.Notes),
		Attachments:  []Attachment{attach},
	}

	return message, nil
}
