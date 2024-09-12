// filename: cmd/o.go

package cmd

import (
    "strings"

    "github.com/neh-cli/neh/cmd/shared"
    "github.com/spf13/cobra"
)

func init() {
    cmdName := shared.GetCommandName()

    getMessage := func(args []string) string {
        return strings.Join(args, " ")
    }

    var dynamicCmd = &cobra.Command{
        Use:   cmdName,
        Short: "Send an inquiry message to the AI",
        Run:   shared.RunDynamicCmd(cmdName, getMessage),
    }
    rootCmd.AddCommand(dynamicCmd)
}
