package cmd

import (
	"fmt"
	"os"

	"github.com/james-do2024/ghi/client"
	"github.com/james-do2024/ghi/config"
	"github.com/james-do2024/ghi/tui"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view [owner/repo] [path]",
	Short: "Non-interactively poll and display GitHub repository content",
	Long: `The view command displays the requested contents simply and then exits,
foregoing syntax highlighting or any further interactive use.`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var sv ViewFunction
		sv = simpleView

		if len(args) < 1 {
			cmd.Help()
			os.Exit(config.ExitErr)
		}
		cmdMain(args, sv)
	},
	DisableFlagsInUseLine: true, // Command takes no flags beyond 'help'
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func simpleView(ts *tui.TuiState, req *client.RestRequest) {
	if ts.FileContent != nil {
		fmt.Print(*ts.FileContent)
	} else {
		fmt.Printf("path: %s\n\n", req.CurrentPath)
		for _, entry := range ts.DirMap {
			fmt.Println(entry)
		}
	}

}
