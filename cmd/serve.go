// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package cmd

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ljvmiranda921/burnout-barometer/pkg"
	"github.com/spf13/cobra"
)

// ServeCommand starts the server that handles request coming from Slack.
func ServeCommand() *cobra.Command {

	var (
		port    int
		cfgPath string
	)

	var command = &cobra.Command{
		Use:     "serve",
		Short:   "Start the server",
		Example: "barometer serve --port=8080",
		RunE: func(cmd *cobra.Command, args []string) error {
			initLogger(verbosity)

			server := pkg.Server{
				Port:    port,
				Router:  httprouter.New(),
				CfgPath: cfgPath,
			}

			server.Routes()
			server.Start()
			return nil
		},
	}

	// Add flags
	command.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	command.Flags().StringVarP(&cfgPath, "config", "c", "config.json", "Path to configuration file")
	return command
}
