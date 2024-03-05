/*
Copyright © 2024 James Taylor <james.taylor@fastmail.com>
*/
package cmd

import (
	"os"

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
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("debug", "d", false, "Run in debug mode")
}
