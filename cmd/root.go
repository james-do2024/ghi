package cmd

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/james-do2024/ghi/client"
	"github.com/james-do2024/ghi/config"
	"github.com/james-do2024/ghi/tui"
	"github.com/spf13/cobra"
)

type ViewFunction func(*tui.TuiState, *client.RestRequest)

var rootCmd = &cobra.Command{
	Use:   "ghi [explore/view] [org/repo path]",
	Short: "A GitHub CLI application",
	Long: `GitHub Interactive (ghi) is a simple command line application which allows a
user to step through any given GitHub repository interactively.

ghi may also be run non-interactively with its 'view' subcommand, which is
useful in scripting and automation.`,
}

func Execute() {
	// Treat these core errors as fatal every time
	if ttyErr := config.TTYVerify(); ttyErr != nil {
		log.Fatalln(ttyErr)
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func cmdMain(args []string, vf ViewFunction) {
	var owner, repo, path string
	var parseErr, initErr, contentUpdateErr error

	ts := tui.Init()

	if len(args) == 1 {
		owner, repo, path, parseErr = handleOneArg(args[0])
	} else {
		owner, repo, path, parseErr = handleTwoArgs(args[0], args[1])
	}

	if parseErr != nil {
		log.Fatalln(parseErr)
	}

	c, initErr := client.NewRestClient()
	if initErr != nil {
		log.Fatalf("unable to initialize REST client: %v\n", initErr) // Simply discontinuing here
	}
	req := &client.RestRequest{
		Owner:       owner,
		Repo:        repo,
		CurrentPath: path,
		GitClient:   c,
		Ctx:         context.Background(),
	}

	contentUpdateErr = ts.UpdateContent(req)
	if contentUpdateErr != nil {
		log.Fatalln(contentUpdateErr)
	}

	vf(ts, req)
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
