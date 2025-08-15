// File: cmd/explain.go

package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var (
	explainModel string
)

var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "This command sends the contents of the clipboard to the AI and receives an explanation from the AI.",
	Run:   runExplainCmd,
}

func init() {
	rootCmd.AddCommand(explainCmd)
	explainCmd.Flags().StringVar(&explainModel, "model", "", "Specify the AI model to use (e.g., gpt-4, claude-3)")
}

func runExplainCmd(cmd *cobra.Command, args []string) {
	clipboardMessage, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	queryMessage := ""
	err = shared.ExecuteWebSocketCommand("explain", queryMessage, clipboardMessage, explainModel)

	if err != nil {
		fmt.Println(err)
	}
}
