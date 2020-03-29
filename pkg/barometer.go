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

	"cloud.google.com/go/bigquery"
	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"
)

const (
	defaultMessage = "Thank you for trusting me"
	ackPrefix      = "Gotcha, I logged your mood"
)

// UpdateLog accepts the userID and the text, parses the timestamp, and stores it into the database.
// If debug is true, then inserting into the database is skipped. This is useful for testing.
func UpdateLog(userID, text string, timestamp time.Time, db Database, twitterClient *twitter.Client, debug bool) (*Message, error) {
	m, notes := ParseMessage(text)
	measure, err := strconv.Atoi(m)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("strconv")
		return nil, err
	}

	item := LogItem{
		Timestamp:     timestamp,
		UserID:        userID,
		LogMeasure:    measure,
		Notes:         notes,
		TwitterClient: twitterClient,
	}

	if debug {
		log.Info("DebugOnly is set to true, will not insert to database")
	} else {
		if err := item.Insert(db); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("logItem.insert")
			return nil, err
		}
	}

	return item.Reply()
}

// ParseMessage extracts the barometer measure and notes from a given text.
func ParseMessage(text string) (string, string) {
	list := strings.Fields(text)
	measure := list[0]
	notes := strings.Join(list[1:], " ")
	return measure, notes
}

// LogItem is the user log for the barometer. This also serves as
// the schema for the database.
type LogItem struct {
	Timestamp     time.Time
	UserID        string
	LogMeasure    int
	Notes         string
	TwitterClient *twitter.Client
}

// Save allows us to implement BigQuery's ValueSaver interface.
func (i *LogItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"timestamp":   i.Timestamp,
		"user_id":     i.UserID,
		"log_measure": i.LogMeasure,
		"notes":       i.Notes,
	}, "", nil
}

// Insert puts the item entry into the specified database.
func (i *LogItem) Insert(db Database) error {
	if err := db.Insert(*i); err != nil {
		log.Errorf("error in inserting item: %v", err)
		return err
	}

	return nil
}

// Reply prepares the Slack message as a response to a slash command.
func (i *LogItem) Reply() (*Message, error) {
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
func (i *LogItem) fetchTwitterMessage(screenName string, count int, userOnly bool) string {
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
