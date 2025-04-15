package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var tCmd = &cobra.Command{
	Use:   "t",
	Short: "This command translates the contents of the clipboard into the language you are using.",
	Run:   runTCmd,
}

func init() {
	rootCmd.AddCommand(tCmd)
}

func runTCmd(cmd *cobra.Command, args []string) {
	clipboardMessage, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	queryMessage := ""
	err = shared.ExecuteWebSocketCommand("t", queryMessage, clipboardMessage)

	if err != nil {
		fmt.Println(err)
	}
}
