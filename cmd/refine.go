// File: cmd/refine.go

package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var refineCmd = &cobra.Command{
	Use:   "refine",
	Short: "This command refines the text in the clipboard.",
	Run:   runRefineCmd,
}

func init() {
	rootCmd.AddCommand(refineCmd)
}

func runRefineCmd(cmd *cobra.Command, args []string) {
	clipboardMessage, err := clipboard.ReadAll()

	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	queryMessage := ""
	err = shared.ExecuteWebSocketCommand("refine", queryMessage, clipboardMessage)

	if err != nil {
		fmt.Println(err)
	}
}
