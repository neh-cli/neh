// File: cmd/fix.go

package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var (
	fixModel string
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "This command sends content that seems to have problems to an LLM (Language Model) and requests proposed solutions.",
	Run:   runFixCmd,
}

func init() {
	rootCmd.AddCommand(fixCmd)
	fixCmd.Flags().StringVar(&fixModel, "model", "", "Specify the AI model to use (e.g., gpt-4, claude-3)")
}

func runFixCmd(cmd *cobra.Command, args []string) {
	clipboardMessage, err := clipboard.ReadAll()

	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	queryMessage := ""
	err = shared.ExecuteWebSocketCommand("fix", queryMessage, clipboardMessage, fixModel)

	if err != nil {
		fmt.Println(err)
	}
}
