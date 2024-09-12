// filename: cmd/decache.go

package cmd

import (
    "fmt"

    "github.com/neh-cli/neh/cmd/shared"
    "github.com/spf13/cobra"
)

var decacheCmd = &cobra.Command{
    Use:   "decache",
    Short: "Remove all query history",
    Long: `This command deletes all previous query history stored on the server.
Use this command to clear any saved interactions with the AI.`,
    Run: runDecacheCmd,
}

func init() {
    rootCmd.AddCommand(decacheCmd)
}

func runDecacheCmd(cmd *cobra.Command, args []string) {
    err := shared.ExecuteWebSocketCommand("decache", "", false)

    if err != nil {
        fmt.Println(err)
    }
}
