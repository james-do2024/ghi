/*
Copyright © 2024 James Taylor <james.taylor@fastmail.com>
*/
package cmd

import (
	"log"
	"os"

	"github.com/james-do2024/ghi/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghi",
	Short: "A GitHub CLI application",
	Long: `GitHub Interactive (ghi) is a simple command line application which allows a
user to step through any given GitHub repository interactively.

ghi may also be run non-interactively with its 'view' subcommand, which is
useful in scripting and automation.`,
}

func Execute() {
	if ttyErr := config.TTYVerify(); ttyErr != nil {
		log.Fatalln(ttyErr)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(config.ExitErr)
	}
}

func init() {
	rootCmd.Flags().BoolP("debug", "d", false, "Run in debug mode")
}
