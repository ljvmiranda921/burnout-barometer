// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"
)

// Configuration contains all important settings for running the command.
type Configuration struct {
	Table string `json:"TABLE"`       // Database URL
	Token string `json:"SLACK_TOKEN"` // Slack token provided by the app for verification
	Area  string `json:"AREA"`        // IANA-compliant area

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
		log.WithFields(log.Fields{"err": err}).Error("os.Create")
		return err
	}
	defer file.Close()

	e := json.NewEncoder(file)
	if err := e.Encode(cfg); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("json.NewEncoder.Encode")
		return err
	}
	return nil
}

func (cfg *Configuration) update(field string, value string) {
	v := reflect.ValueOf(cfg).Elem().FieldByName(field)
	if v.IsValid() {
		v.SetString(value)
	}
}

// ReadConfiguration reads the configuration file and returns an instance
// of a Configuration.
func ReadConfiguration(cfgPath string) (*Configuration, error) {

	// Open configuration file
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("os.Open")
		return nil, err
	}

	// Create a decoder from the file
	d := json.NewDecoder(cfgFile)
	config := &Configuration{}
	if err = d.Decode(config); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("json.NewDecoder.Decode")
		return nil, err
	}

	// Decode secrets
	secretFields := []string{
		"Token",
		"TwitterConsumerKey",
		"TwitterConsumerSecret",
		"TwitterAccessKey",
		"TwitterAccessSecret",
	}

	for _, field := range secretFields {
		v := reflect.ValueOf(*config)
		enc := v.FieldByName(field).String()
		dec, err := base64.StdEncoding.DecodeString(enc)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("base64.StdEncoding.DecodeString")
			return nil, err
		}

		config.update(field, string(dec))
	}
	return config, nil
}
