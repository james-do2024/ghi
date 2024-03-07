/*
Copyright © 2024 James Taylor <james.taylor@fastmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Non-interactively poll and display GitHub repository content",
	Long: `The view command displays the requested contents simply and then exits,
foregoing syntax highlighting or any further interactive use.`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var sv ViewFunction
		sv = simpleView

		cmdMain(args, sv)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func simpleView(file *string, dir []*github.RepositoryContent) {
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

}
