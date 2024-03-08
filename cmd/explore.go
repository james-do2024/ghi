/*
Copyright © 2024 James Taylor <james.taylor@fastmail.com>
*/
package cmd

import (
	"log"
	"os"

	"github.com/james-do2024/ghi/client"
	"github.com/james-do2024/ghi/config"
	"github.com/james-do2024/ghi/tui"
	"github.com/spf13/cobra"
)

var exploreCmd = &cobra.Command{
	Use:   "explore [owner/repo] [path]",
	Short: "Interactively poll and display GitHub repository content",
	Long: `The 'explore' command fetches and displays repository content, allowing
the user to navigate the repository with numbers and symbols.

If the content is a file, it displays with syntax highlighting and paged output is
attempted.

Keystrokes are as follows:
[0-9] : Select entry by number (can also be 10 or higher)
^  : The caret character navigates to the root of the repository
.. : Navigates one directory back, if possible
q  : Quits the interactive session

For non-interactive use, the 'view' command is what you want.`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var rv ViewFunction
		rv = richView

		if len(args) < 1 {
			cmd.Help()
			os.Exit(config.ExitErr)
		}
		cmdMain(args, rv)
	},
	DisableFlagsInUseLine: true, // Command takes no flags beyond 'help'
}

func init() {
	rootCmd.AddCommand(exploreCmd)
}

func richView(ts *tui.TuiState, req *client.RestRequest) {
	err := ts.Interact(req)
	if err != nil {
		log.Fatalln(err)
	}
}
