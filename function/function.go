// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

// Package function contains a single function that can be deployed as a
// serverless microservice.
package function

import (
	"encoding/json"
	"net/http"

	"github.com/ljvmiranda921/burnout-barometer/pkg"
	log "github.com/sirupsen/logrus"
)

// BurnoutBarometerFn is a Cloud Function that takes a log message from a Slack
// slash command and stores it into BigQuery as a streaming insert.
func BurnoutBarometerFn(w http.ResponseWriter, r *http.Request) {
	log.Info("request received")

	// Setup application variables
	config, err := pkg.ReadConfiguration("config.json")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("ReadConfiguration")
	}

	// Validate request and parse the submitted form
	if r.Method != "POST" {
		http.Error(w, "only POST requests are accepted", 405)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "couldn't parse form", 400)
		log.WithFields(log.Fields{"err": err}).Fatal("http.Request.ParseForm")
	}

	if err := pkg.VerifyWebhook(r.Form, config.Token); err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("verifyWebHook")
	}

	if len(r.Form["text"]) == 0 {
		log.Fatal("empty text in form")
	}

	// Store the message and timestamp to BigQuery
	db, err := pkg.NewDatabase(config.Table)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("NewDatabase")
	}
	req := &pkg.Request{
		Text:      r.Form["text"][0],
		UserID:    r.Form["user_id"][0],
		Timestamp: r.Header.Get("X-Slack-Request-Timestamp"),
		Area:      config.Area,
		DB:        db,
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
