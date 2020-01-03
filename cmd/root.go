// Copyright © 2019 Lester James V. Miranda <ljvmiranda@gmail.com>
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

// Package cmd contains all commands in the burnout-barometer CLI
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var verbosity int
var config string

// NewCommand returns a new instance of an optserve command.
func NewCommand() *cobra.Command {

	var command = &cobra.Command{
		Use:   "barometer [command]",
		Short: "barometer is the command-line interface to Burnout Barometer",
		Run: func(cmd *cobra.Command, args []string) {
			initLogger(verbosity)
			cmd.HelpFunc()(cmd, args)
		},
	}

	// Define persistent flags
	command.PersistentFlags().CountVarP(&verbosity, "verbosity", "v", "set verbosity")

	// Add subcommands
	command.AddCommand(ServeCommand())

	return command
}

func initLogger(verbosity int) {
	switch {
	case verbosity == 1:
		log.SetLevel(log.DebugLevel)
	case verbosity >= 2:
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.InfoLevel)

	}
}
