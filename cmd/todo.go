// File: cmd/todo.go

package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "This command sends text containing TODO to an AI and asks for proposed solutions to the TODO.",
	Run:   runTodoCmd,
}

func init() {
	rootCmd.AddCommand(todoCmd)
}

func runTodoCmd(cmd *cobra.Command, args []string) {
	clipboardMessage, err := clipboard.ReadAll()

	if err != nil {
		fmt.Println("Operation aborted because clipboard contents could not be retrieved.", err)
		return
	}

	queryMessage := ""
	err = shared.ExecuteWebSocketCommand("todo", queryMessage, clipboardMessage, "")

	if err != nil {
		fmt.Println(err)
	}
}
