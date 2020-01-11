// Copyright 2019 Lester James V. Miranda. All rights reserved.
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
	ProjectID string `json:"PROJECT_ID"`
	Table     string `json:"BQ_TABLE"`
	Token     string `json:"SLACK_TOKEN"`
	Area      string `json:"AREA"`
}

// WriteConfiguration creates a configuration file at a given output path.
func (cfg *Configuration) WriteConfiguration(outputPath string) error {

	file, _ := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
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
