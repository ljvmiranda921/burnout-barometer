package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
)

type configuration struct {
	ProjectID string `json:"PROJECT_ID"`
	Table     string `json:"BQ_TABLE"`
	Token     string `json:"SLACK_TOKEN"`
}

func setup(ctx context.Context) {
	if config == nil {
		cfgFile, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("os.Open: %v", err)
		}

		d := json.NewDecoder(cfgFile)
		config = &configuration{}
		if err = d.Decode(config); err != nil {
			log.Fatalf("Decode: %v", err)
		}
	}
}
