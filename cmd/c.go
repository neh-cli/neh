package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var cCmd = &cobra.Command{
	Use:   "c",
	Short: "Send the contents of the clipboard and a query about that content to the LLM",
	Run:   runCCmd,
}

func init() {
	rootCmd.AddCommand(cCmd)
}

func runCCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a message to send")
		return
	}

	queryMessage := strings.Join(args, " ")
	if os.Getenv("NEH_DEBUG") == "t" {
		fmt.Printf("Query Message: %s\n", queryMessage)
	}

	// Retrieve the contents of the clipboard and store them in clipboardMessage
	clipboardMessage, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	err = shared.ExecuteWebSocketCommand("c", queryMessage, clipboardMessage)

	if err != nil {
		fmt.Println(err)
	}
}
