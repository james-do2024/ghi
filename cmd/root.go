/*
Copyright © 2024 James Taylor <james.taylor@fastmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

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

func handleOneArg(arg string) (owner, repo, path string, err error) {
	// Allow users to enter: ghi <view, explore> owner/repo [path]
	baseRegex := regexp.MustCompile(`^[\w-]+/[\w-]+$`)

	// Another pattern is: ghi <view, explore> owner/repo/path
	oneArgRegex := regexp.MustCompile(`^[\w-]+/[\w-]+(?:/[\w-./]+)$`)

	if oneArgRegex.MatchString(arg) {
		splitArg := strings.SplitN(arg, "/", 3)
		owner, repo, path = splitArg[0], splitArg[1], splitArg[2]
		return owner, repo, path, nil
	} else if baseRegex.MatchString(arg) {
		splitArg := strings.SplitN(arg, "/", 2)
		owner, repo = splitArg[0], splitArg[1]
		return owner, repo, "", nil
	}
	return "", "", "", fmt.Errorf("Input does not match `owner/repo` pattern: %s", arg)
}

func handleTwoArgs(arg1, arg2 string) (owner, repo, path string, err error) {
	baseRegex := regexp.MustCompile(`^[\w-]+/[\w-]+$`)
	if baseRegex.MatchString(arg1) {
		splitArg := strings.SplitN(arg1, "/", 2)
		owner, repo = splitArg[0], splitArg[1]
		path = arg2
		return owner, repo, path, nil
	}
	return "", "", "", fmt.Errorf("Input does not match `owner/repo` pattern: %s", arg1)
}
