// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package function

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ljvmiranda921/burnout-barometer/pkg"
	log "github.com/sirupsen/logrus"
)

// BurnoutBarometerFn is a Cloud Function that takes a log message from a Slack
// slash command and stores it into BigQuery as a streaming insert.
func BurnoutBarometerFn(w http.ResponseWriter, r *http.Request) {
	log.Info("request received")

	// Setup application variables
	config, err := pkg.NewConfiguration(r.Context(), "config.json")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("NewConfiguration")
	}

	// Validate request and parse the submitted form
	if r.Method != "POST" {
		http.Error(w, "only POST requests are accepted", 405)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "couldn't parse form", 400)
		log.WithFields(log.Fields{"err": err}).Fatal("http.Request.ParseForm")
	}

	log.Infof("project-id: %s", config.ProjectID)
	log.Infof("table: %s", config.Table)
	log.Infof("area: %s", config.Area)
	log.Infof("token: %s", config.Token)

	if err := verifyWebHook(r.Form, config.Token); err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("verifyWebHook")
	}

	if len(r.Form["text"]) == 0 {
		log.Fatal("empty text in form")
	}

	// Store the message and timestamp to BigQuery
	req := &pkg.Request{
		Text:      r.Form["text"][0],
		UserID:    r.Form["user_id"][0],
		Timestamp: r.Header.Get("X-Slack-Request-Timestamp"),
		Area:      config.Area,
		BQTable:   config.Table,
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

func verifyWebHook(form url.Values, token string) error {
	t := form.Get("token")
	if len(t) == 0 {
		return fmt.Errorf("empty form token")
	}

	if t != token {
		return fmt.Errorf("invalid request/credentials: %q", t)
	}

	return nil
}
