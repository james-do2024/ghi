/*
Copyright © 2024 James Taylor <james.taylor@fastmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/james-do2024/ghi/client"
	"github.com/james-do2024/ghi/config"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Non-interactively poll and display GitHub repository content",
	Long: `The view command displays the requested contents simply and then exits,
foregoing syntax highlighting or any further interactive use.`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cmdMain(args)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func cmdMain(args []string) {
	var owner, repo, path string
	var err error

	if len(args) == 1 {
		owner, repo, path, err = handleOneArg(args[0])
	} else {
		owner, repo, path, err = handleTwoArgs(args[0], args[1])
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(config.ExitErr)
	}
	runRequest(owner, repo, path)
}

func runRequest(owner, repo, path string) error {
	ctx := context.Background()
	c, err := client.NewRestClient()
	if err != nil {
		return err
	}

	req := &client.RestRequest{
		GitClient: c,
		Ctx:       ctx,
	}

	file, dir, err := req.GetContent(owner, repo, path)
	if err != nil {
		return err
	}

	if file != nil {
		fmt.Print(*file)
	} else {
		for _, entry := range dir {
			if *entry.Type == "dir" {
				*entry.Name += "/" // Simple way to differentiate directories in repo view
			}
			fmt.Println(*entry.Name)
		}
	}
	return nil
}
