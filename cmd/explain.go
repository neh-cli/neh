// File: cmd/explain.go

package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "This command sends the contents of the clipboard to the AI and receives an explanation from the AI.",
	Run:   runExplainCmd,
}

func init() {
	rootCmd.AddCommand(explainCmd)
}

func runExplainCmd(cmd *cobra.Command, args []string) {
	// Retrieve the contents of the clipboard and store them in clipboardMessage
	clipboardMessage, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	queryMessage := ""
	err = shared.ExecuteWebSocketCommand("explain", queryMessage, clipboardMessage)

	if err != nil {
		fmt.Println(err)
	}
}
