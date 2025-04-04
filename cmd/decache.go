// File: cmd/decache.go
package cmd

import (
	"fmt"

	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

// The decacheCmd is a command to remove all query history.
// By using this command, you can delete all previous query history
// stored on the server.
var decacheCmd = &cobra.Command{
	Use:   "decache",
	Short: "Remove all query history",
	Long: `This command deletes all previous query history stored on the server.
Use this command to clear any saved interactions with the AI.`,
	Run: runDecacheCmd,
}

// The init function is used to add the decacheCmd to the rootCmd.
// By using this function, the decache command will be recognized
// as part of the CLI.
func init() {
	rootCmd.AddCommand(decacheCmd)
}

// The runDecacheCmd function is called when the decache command is executed.
// This function uses the shared.ExecuteWebSocketCommand function to send
// the "decache" command to the server.
// It is set to send the command without waiting for a response from the server.
func runDecacheCmd(cmd *cobra.Command, args []string) {
	queryMessage := ""
	err := shared.ExecuteWebSocketCommand("decache", queryMessage)

	if err != nil {
		fmt.Println(err)
	}
}
