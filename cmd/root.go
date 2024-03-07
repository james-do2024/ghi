/*
Copyright © 2024 James Taylor <james.taylor@fastmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v60/github"
	"github.com/james-do2024/ghi/client"
	"github.com/james-do2024/ghi/config"
	"github.com/spf13/cobra"
)

type ViewFunction func(*string, []*github.RepositoryContent)

var rootCmd = &cobra.Command{
	Use:   "ghi",
	Short: "A GitHub CLI application",
	Long: `GitHub Interactive (ghi) is a simple command line application which allows a
user to step through any given GitHub repository interactively.

ghi may also be run non-interactively with its 'view' subcommand, which is
useful in scripting and automation.`,
}

func init() {
	rootCmd.Flags().BoolP("debug", "d", false, "Run in debug mode")
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

func cmdMain(args []string, vf ViewFunction) {
	var owner, repo, path string
	var file *string
	var dir []*github.RepositoryContent
	var err error

	if len(args) == 1 {
		owner, repo, path, err = handleOneArg(args[0])
	} else {
		owner, repo, path, err = handleTwoArgs(args[0], args[1])
	}

	if err == nil {
		file, dir, err = runRequest(owner, repo, path)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(config.ExitErr)
	}

	vf(file, dir)
}

func runRequest(owner, repo, path string) (*string, []*github.RepositoryContent, error) {
	ctx := context.Background()
	c, err := client.NewRestClient()
	if err != nil {
		return nil, nil, err
	}

	req := &client.RestRequest{
		GitClient: c,
		Ctx:       ctx,
	}

	file, dir, err := req.GetContent(owner, repo, path)
	if err != nil {
		return nil, nil, err
	}

	return file, dir, nil
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
