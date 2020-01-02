package barometer

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
	Area      string `json:"AREA"`
}

var config *configuration

func setup(ctx context.Context) error {
	if config == nil {
		cfgFile, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("os.Open: %v", err)
			return err
		}

		d := json.NewDecoder(cfgFile)
		config = &configuration{}
		if err = d.Decode(config); err != nil {
			log.Fatalf("Decode: %v", err)
			return err
		}
	}

	return nil
}
