// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ljvmiranda921/burnout-barometer/pkg"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// InitCommand creates a configuration file from a prompt or through environment variables
func InitCommand() *cobra.Command {

	var (
		useEnvVars bool
		outputPath string
	)

	var command = &cobra.Command{
		Use:   "init",
		Short: "Initialize a configuration file",
		Long: `
This command creates a configuration file, config.json, that will be used later
on when running the server. The values can be obtained from the prompt or
through environment variables (prefixed with 'BB_*').

If you're planning to use environment variables, then use the --use-env-vars
flag. Be sure that you have set the following:

- BB_PROJECT_ID: the Google Cloud Platform project ID
- BQ_TABLE: the BigQuery table to store the logs
- BB_SLACK_TOKEN: the Slack Token for verifying requests
- BB_AREA: the local area for setting the timestamp
`,
		Example: "barometer init",
		RunE: func(cmd *cobra.Command, args []string) error {
			initLogger(verbosity)

			if useEnvVars {
				cfg, err := getConfigFromEnvs()
				if err != nil {
					return err
				}

				if err := pkg.WriteConfiguration(cfg, outputPath); err != nil {
					return err
				}
				fmt.Println("Configuration file generated in %s!", outputPath)

			} else {
				cfg, err := getConfigFromPrompt()
				if err != nil {
					return err
				}

				if err := pkg.WriteConfiguration(cfg, outputPath); err != nil {
					return err
				}
				fmt.Printf("Configuration file generated in %s!", outputPath)
			}

			return nil
		},
	}

	// Add flags
	command.Flags().BoolVar(&useEnvVars, "use-env-vars", false, "Use environment variables")
	command.Flags().StringVarP(&outputPath, "output-path", "o", "config.json", "Output path for writing configuration file")
	return command
}

func getConfigFromEnvs() (*pkg.Configuration, error) {

	projectID, err := lookupEnvVar("BB_PROJECT_ID")
	if err != nil {
		return nil, err
	}

	table, err := lookupEnvVar("BB_TABLE")
	if err != nil {
		return nil, err
	}

	slackToken, err := lookupEnvVar("BB_SLACK_TOKEN")
	if err != nil {
		return nil, err
	}
	slackTokenEnc := base64.StdEncoding.EncodeToString([]byte(slackToken))

	area, err := lookupEnvVar("BB_AREA")
	if err != nil {
		return nil, err
	}

	config := &pkg.Configuration{
		ProjectID: projectID,
		Table:     table,
		Token:     slackTokenEnc,
		Area:      area,
	}

	return config, nil
}

func lookupEnvVar(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}

	err := fmt.Errorf("Cannot find environment variable: %s", key)
	return "", err
}

func getConfigFromPrompt() (*pkg.Configuration, error) {
	projectID, err := promptString("GCP Project ID", "my-gcp-project", false)
	if err != nil {
		return nil, err
	}

	table, err := promptString("Table to store all logs", "bq://my-gcp-project.my-dataset.my-table", false)
	if err != nil {
		return nil, err
	}

	slackToken, err := promptString("Slack verification token", "", true)
	if err != nil {
		return nil, err
	}
	slackTokenEnc := base64.StdEncoding.EncodeToString([]byte(slackToken))

	area, err := promptString("Where are you? (Refer to IANA Timezone database)", "Asia/Manila", false)
	if err != nil {
		return nil, err
	}

	config := &pkg.Configuration{
		ProjectID: projectID,
		Table:     table,
		Token:     slackTokenEnc,
		Area:      area,
	}
	return config, nil
}

func promptString(name string, defaultVal string, mask bool) (string, error) {

	var prompt promptui.Prompt

	validate := func(input string) error {
		if len(strings.TrimSpace(input)) < 1 {
			return errors.New("Input must not be empty")
		}
		return nil
	}

	if mask {
		prompt = promptui.Prompt{
			Label:    name,
			Default:  defaultVal,
			Validate: validate,
			Mask:     '*',
		}
	} else {
		prompt = promptui.Prompt{
			Label:    name,
			Default:  defaultVal,
			Validate: validate,
		}
	}

	return prompt.Run()
}
