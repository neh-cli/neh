// File: cmd/refactor.go

package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var (
	refactorModel string
)

var refactorCmd = &cobra.Command{
	Use:   "refactor",
	Short: "Send the source code in the clipboard to LLM and request refactoring.",
	Run:   runRefactorCmd,
}

func init() {
	rootCmd.AddCommand(refactorCmd)
	refactorCmd.Flags().StringVar(&refactorModel, "model", "", "Specify the AI model to use (e.g., gpt-4.1, gpt-5)")
}

func runRefactorCmd(cmd *cobra.Command, args []string) {
	clipboardMessage, err := clipboard.ReadAll()

	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	queryMessage := ""
	err = shared.ExecuteWebSocketCommand("refactor", queryMessage, clipboardMessage, refactorModel)

	if err != nil {
		fmt.Println(err)
	}
}
