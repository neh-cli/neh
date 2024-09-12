// filename: cmd/decache.go

package cmd

import (
    "github.com/neh-cli/neh/cmd/shared"
    "github.com/spf13/cobra"
)

func init() {
    cmdName := shared.GetCommandName()

    getMessage := func(args []string) string {
        return ""
    }

    var dynamicCmd = &cobra.Command{
        Use:   cmdName,
        Short: "Remove all query history",
        Long: `This command deletes all previous query history stored on the server.
-Use this command to clear any saved interactions with the AI.`,
        Run:   shared.RunDynamicCmd(cmdName, getMessage),
    }
    rootCmd.AddCommand(dynamicCmd)
}
