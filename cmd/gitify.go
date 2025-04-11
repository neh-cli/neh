// File: cmd/gitify.go

package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var gitifyCmd = &cobra.Command{
	Use:   "gitify",
	Short: "This command suggests the title of the issue you are currently working on, the git branch name, and even the git commit message from the contents of the clipboard.",
	Run:   runGitifyCmd,
}

func init() {
	rootCmd.AddCommand(gitifyCmd)
}

func runGitifyCmd(cmd *cobra.Command, args []string) {
	clipboardMessage, err := clipboard.ReadAll()

	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	queryMessage := ""
	err = shared.ExecuteWebSocketCommand("gitify", queryMessage, clipboardMessage)

	if err != nil {
		fmt.Println(err)
	}
}
