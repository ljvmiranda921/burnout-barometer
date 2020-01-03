// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

// Package pkg contains types and methods for interacting with the barometer.
package pkg

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
)

// Request defines the common form parameters when a slash command is invoked.
type Request struct {
	Text      string
	UserID    string
	Timestamp string
	Area      string
	BQTable   string
	Item      Log
}

// Process parses the request and stores to BigQuery.
func (r *Request) Process() (*Message, error) {
	m, notes, err := r.ParseMessage()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("Request.ParseMessage")
		return nil, err
	}

	ts, err := r.GetTimestamp()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("Request.GetTimestamp")
		return nil, err
	}

	measure, err := strconv.Atoi(m)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("strconv")
		return nil, err
	}

	r.Item = Log{
		Timestamp:  ts,
		UserID:     r.UserID,
		LogMeasure: measure,
		Notes:      notes,
	}

	if err := r.InsertToTable(); err != nil {
		log.Fatalf("error in InsertToTable: %v", err)
		return nil, err
	}

	return r.Item.FormatReply()
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

// InsertToTable adds the Item entry into the specified Bigquery table.
func (r *Request) InsertToTable() error {
	ctx := context.Background()
	projectID, datasetID, tableID := r.splitBQPath(r.BQTable)
	log.Printf("using BQ table: %s", r.BQTable)
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("error in bigquery.NewClient: %v", err)
	}

	inserter := client.Dataset(datasetID).Table(tableID).Inserter()
	items := []*Log{&r.Item}

	if err := inserter.Put(ctx, items); err != nil {
		return err
	}
	return nil
}

func (r *Request) splitBQPath(p string) (string, string, string) {
	s := strings.Split(p, ".")
	return s[0], s[1], s[2]
}

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
// see https://api.slack.com/docs/message-formatting for more information.
type Message struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

// Attachment defines the message output after running the slash command.
type Attachment struct {
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
}
