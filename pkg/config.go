// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

// Configuration contains all important settings for running the command
type Configuration struct {
	ProjectID string `json:"PROJECT_ID"`
	Table     string `json:"BQ_TABLE"`
	Token     string `json:"SLACK_TOKEN"`
	Area      string `json:"AREA"`
}

// Setup reads the configuration file and updates the parameters accordingly
func (c *Configuration) Setup(ctx context.Context) error {
	if c == nil {
		cfgFile, err := os.Open("config.json")
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Fatal("os.Open")
			return err
		}

		d := json.NewDecoder(cfgFile)
		if err = d.Decode(c); err != nil {
			log.WithFields(log.Fields{"err": err}).Fatal("json.NewDecoder.Decode")
			return err
		}

		// Decode `SLACK_TOKEN`
		slackToken, err := base64.StdEncoding.DecodeString(c.Token)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Fatal("base64.StdEncoding.DecodeString")
			return err
		}
		c.Token = string(slackToken)
	}

	return nil
}
