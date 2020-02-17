// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package cmd

import (
	"encoding/base64"
	"encoding/json"
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
flag. Find all available options in this link: https://ljvmiranda921.github.io/burnout-barometer/installation/
`,
		Example: "barometer init",
		RunE: func(cmd *cobra.Command, args []string) error {
			initLogger(verbosity)

			if useEnvVars {
				cfg, err := getConfig(fromEnvVar)
				if err != nil {
					return err
				}

				if err := cfg.WriteConfiguration(outputPath); err != nil {
					return err
				}
				fmt.Printf("Configuration file generated in %s!", outputPath)

			} else {
				cfg, err := getConfig(fromPrompt)
				if err != nil {
					return err
				}

				if err := cfg.WriteConfiguration(outputPath); err != nil {
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

// Specify all configuration options here. As of now, we can only maintain this
// by manually checking the struct fields within pkg.Configuration. It may be
// easier to decouple them instead.
type opt struct {
	name     string
	toEncode bool

	// environment-var specific options
	envVarName string

	// promptui specific options
	prompt     string
	defaultVal string
	mask       bool
}

var opts = []opt{
	// TODO: Remove PROJECT_ID
	opt{
		name:       "PROJECT_ID",
		toEncode:   false,
		envVarName: "BB_PROJECT_ID",
		prompt:     "GCP Project ID",
		defaultVal: "my-project-id",
		mask:       false,
	},
	opt{
		name:       "TABLE",
		toEncode:   false,
		envVarName: "BB_TABLE",
		prompt:     "Database URL to store all logs",
		defaultVal: "bq://my-gcp-project.my-dataset.my-table",
		mask:       false,
	},
	opt{
		name:       "SLACK_TOKEN",
		toEncode:   true,
		envVarName: "BB_SLACK_TOKEN",
		prompt:     "Slack verification token",
		defaultVal: "",
		mask:       true,
	},
	opt{
		name:       "AREA",
		toEncode:   false,
		envVarName: "BB_AREA",
		prompt:     "Where are you? (refer to IANA timezone database)",
		defaultVal: "Asia/Manila",
		mask:       false,
	},
	opt{
		name:       "TWITTER_CONSUMER_KEY",
		toEncode:   true,
		envVarName: "BB_TWITTER_CONSUMER_KEY",
		prompt:     "Twitter API Consumer Key",
		defaultVal: "",
		mask:       true,
	},
	opt{
		name:       "TWITTER_CONSUMER_SECRET",
		toEncode:   true,
		envVarName: "BB_TWITTER_CONSUMER_SECRET",
		prompt:     "Twitter API Consumer Secret",
		defaultVal: "",
		mask:       true,
	},
	opt{
		name:       "TWITTER_ACCESS_KEY",
		toEncode:   true,
		envVarName: "BB_TWITTER_ACCESS_KEY",
		prompt:     "Twitter API Access Key",
		defaultVal: "",
		mask:       true,
	},
	opt{
		name:       "TWITTER_ACCESS_SECRET",
		toEncode:   true,
		envVarName: "BB_TWITTER_ACCESS_SECRET",
		prompt:     "Twitter API Access Secret",
		defaultVal: "",
		mask:       true,
	},
}

func getConfig(fn func(options opt) (string, error)) (*pkg.Configuration, error) {
	m := make(map[string]interface{})
	for i := 0; i < len(opts); i++ {
		val, err := fn(opts[i])
		if err != nil {
			return nil, err
		}
		m[opts[i].name] = val
	}

	jsonString, _ := json.Marshal(m) // convert map to json
	config := &pkg.Configuration{}
	json.Unmarshal(jsonString, config) // convert json to struct
	return config, nil
}

func fromEnvVar(options opt) (string, error) {
	if value, exists := os.LookupEnv(options.envVarName); exists {
		if options.toEncode {
			valueEnc := base64.StdEncoding.EncodeToString([]byte(value))
			return valueEnc, nil
		}
		return value, nil
	}

	err := fmt.Errorf("Cannot find environment variable: %s", options.envVarName)
	return "", err
}

func fromPrompt(options opt) (string, error) {

	var prompt promptui.Prompt

	validate := func(input string) error {
		if len(strings.TrimSpace(input)) < 1 {
			return errors.New("Input must not be empty")
		}
		return nil
	}

	if options.mask {
		prompt = promptui.Prompt{
			Label:    options.prompt,
			Default:  options.defaultVal,
			Validate: validate,
			Mask:     '*',
		}
	} else {
		prompt = promptui.Prompt{
			Label:    options.prompt,
			Default:  options.defaultVal,
			Validate: validate,
		}
	}

	value, err := prompt.Run()
	if err != nil {
		return "", nil
	}

	if options.toEncode {
		valueEnc := base64.StdEncoding.EncodeToString([]byte(value))
		return valueEnc, nil
	}

	return value, nil
}
