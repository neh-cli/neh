// File: cmd/o.go
//
// This file contains the implementation of the "o" command, which sends an inquiry message to the AI.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var oCmd = &cobra.Command{
	Use:   "o",
	Short: "Send an inquiry message to the AI",
	Run:   runOCmd,
}

func init() {
	rootCmd.AddCommand(oCmd)
}

func runOCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a message to send")
		return
	}

	queryMessage := strings.Join(args, " ")
	if os.Getenv("NEH_DEBUG") == "t" {
		fmt.Printf("Query Message: %s\n", queryMessage)
	}

	err := shared.ExecuteWebSocketCommand("o", queryMessage, "")

	if err != nil {
		fmt.Println(err)
	}
}
