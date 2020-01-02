package barometer

import (
	"context"
	"encoding/base64"
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
			log.Fatalf("error in os.Open: %v", err)
			return err
		}

		d := json.NewDecoder(cfgFile)
		config = &configuration{}
		if err = d.Decode(config); err != nil {
			log.Fatalf("error in Decode: %v", err)
			return err
		}

		// Decode `SLACK_TOKEN`
		slackToken, err := base64.StdEncoding.DecodeString(config.Token)
		config.Token = string(slackToken)
		if err != nil {
			log.Fatalf("error in DecodeString: %v", err)
		}
	}

	return nil
}
