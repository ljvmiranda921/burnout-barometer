// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

// Package pkg contains types and methods for interacting with the barometer.
package pkg

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"4d63.com/tz"
	"cloud.google.com/go/bigquery"
	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"
)

const (
	defaultMessage = "Thank you for trusting me"
	ackPrefix      = "Gotcha, I logged your mood"
)

// Request defines the common form parameters when a slash command is invoked.
type Request struct {
	Text          string          // The submitted text in the slash command
	UserID        string          // Slack User ID that submitted the request
	Timestamp     string          // Timestamp of the request
	Area          string          // IANA-compliant area
	DB            Database        // Database to insert into
	TwitterClient *twitter.Client // Twitter client

	// If true, then message will not insert into the database. Useful for testing.
	Debug bool
}

// Process parses the request and stores into the Database.
func (r *Request) Process() (*Message, error) {
	m, notes := r.message()

	ts, err := r.timestamp()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Request.getTimestamp")
		return nil, err
	}

	measure, err := strconv.Atoi(m)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("strconv")
		return nil, err
	}

	item := logItem{
		Timestamp:     ts,
		UserID:        r.UserID,
		LogMeasure:    measure,
		Notes:         notes,
		TwitterClient: r.TwitterClient,
	}

	if r.Debug {
		log.Info("DebugOnly is set to true, will not insert to database")
	} else {
		if err := item.insert(r.DB); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("logItem.insert")
			return nil, err
		}
	}

	return item.formatReply()
}

// message extracts the barometer measure and notes from the form text.
func (r *Request) message() (string, string) {
	list := strings.Fields(r.Text)
	measure := list[0]
	notes := strings.Join(list[1:], " ")
	return measure, notes
}

// timestamp obtains the timestamp value from the request.
func (r *Request) timestamp() (time.Time, error) {
	i, err := strconv.ParseInt(r.Timestamp, 10, 64)
	if err != nil {
		log.Errorf("cannot parse timestamp %s: %v", r.Timestamp, err)
		return time.Time{}, err
	}
	loc, err := tz.LoadLocation(r.Area)
	if err != nil {
		log.Errorf("cannot find location: %s", r.Area)
		return time.Time{}, err
	}

	return time.Unix(i, 0).In(loc), nil
}

// logItem is the user log for the barometer. This also serves as
// the schema for the database.
type logItem struct {
	Timestamp     time.Time
	UserID        string
	LogMeasure    int
	Notes         string
	TwitterClient *twitter.Client
}

// Save allows us to implement BigQuery's ValueSaver interface.
func (i *logItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"timestamp":   i.Timestamp,
		"user_id":     i.UserID,
		"log_measure": i.LogMeasure,
		"notes":       i.Notes,
	}, "", nil
}

// insert puts the item entry into the specified database.
func (i *logItem) insert(db Database) error {
	if err := db.Insert(*i); err != nil {
		log.Errorf("error in inserting item: %v", err)
		return err
	}

	return nil
}

// formatReply prepares the Slack message as a response to a slash command.
func (i *logItem) formatReply() (*Message, error) {
	var text string
	if i.TwitterClient != nil {
		text = i.fetchTwitterMessage("tinycarebot", 20, true)
	} else {
		text = defaultMessage
	}

	attach := Attachment{
		Color: "#ef4631",
		Title: "Here's your message from Burnout Barometer",
		Text:  text,
	}

	msg := &Message{
		ResponseType: "ephemeral",
		Text:         fmt.Sprintf("%s: %d (%s)", ackPrefix, i.LogMeasure, i.Notes),
		Attachments:  []Attachment{attach},
	}

	return msg, nil
}

// fetchTwitterMessage gets N number of the latest tweets from a username (preferably, tinycarebot)
func (i *logItem) fetchTwitterMessage(screenName string, count int, userOnly bool) string {
	log.WithFields(log.Fields{"username": screenName}).Trace("fetching tweet")
	tweets, resp, err := i.TwitterClient.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName:     screenName,
		Count:          count,
		ExcludeReplies: &userOnly,
	})

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Tracef("fetch unsuccessful: %v", err)
		return defaultMessage
	}

	// Choose a random tweet from tinycarebot
	rand.Seed(time.Now().Unix())
	tweet := tweets[rand.Intn(len(tweets))]
	log.Tracef("status (%s), tweet: %s", resp.Status, tweet.Text)
	text := fmt.Sprintf("%s (@%s)", tweet.Text, screenName)
	return text
}

// Message is the Slack message event. see
// https://api.slack.com/docs/message-formatting for more information.
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
