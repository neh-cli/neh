// File: cmd/refactor.go

package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var refactorCmd = &cobra.Command{
	Use:   "refactor",
	Short: "Send the source code in the clipboard to LLM and request refactoring.",
	Run:   runRefactorCmd,
}

func init() {
	rootCmd.AddCommand(refactorCmd)
}

func runRefactorCmd(cmd *cobra.Command, args []string) {
	clipboardMessage, err := clipboard.ReadAll()

	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	queryMessage := ""
	err = shared.ExecuteWebSocketCommand("refactor", queryMessage, clipboardMessage)

	if err != nil {
		fmt.Println(err)
	}
}
