// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"encoding/base64"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

// Configuration contains all important settings for running the command.
type Configuration struct {
	ProjectID string `json:"PROJECT_ID"`  // GCP Project ID
	Table     string `json:"TABLE"`       // Database URL
	Token     string `json:"SLACK_TOKEN"` // Slack token provided by the app for verification
	Area      string `json:"AREA"`        // IANA-compliant area

	// This defines the API keys for accessing the Twitter API
	// and get messages from the tiny-care bots
	TwitterConsumerKey    string `json:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret string `json:"TWITTER_CONSUMER_SECRET"`
	TwitterAccessKey      string `json:"TWITTER_ACCESS_KEY"`
	TwitterAccessSecret   string `json:"TWITTER_ACCESS_SECRET"`
}

// WriteConfiguration creates a configuration file at a given output path.
func (cfg *Configuration) WriteConfiguration(outputPath string) error {

	file, err := os.Create(outputPath)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("os.Create")
		return err
	}
	defer file.Close()

	e := json.NewEncoder(file)
	if err := e.Encode(cfg); err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("json.NewEncoder.Encode")
		return err
	}

	return nil
}

// ReadConfiguration reads the configuration file and returns an instance
// of a Configuration.
func ReadConfiguration(cfgPath string) (*Configuration, error) {

	// Open configuration file
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("os.Open")
		return nil, err
	}

	// Create a decoder from the file
	d := json.NewDecoder(cfgFile)
	config := &Configuration{}
	if err = d.Decode(config); err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("json.NewDecoder.Decode")
		return nil, err
	}

	log.Tracef("%v", config)

	// Decode `SLACK_TOKEN`
	slackToken, err := base64.StdEncoding.DecodeString(config.Token)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("base64.StdEncoding.DecodeString")
		return nil, err
	}
	config.Token = string(slackToken)

	return config, nil
}
