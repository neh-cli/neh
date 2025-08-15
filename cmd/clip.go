// File: cmd/clip.go

package cmd

import (
	"fmt"

	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var clipCmd = &cobra.Command{
	Use:   "clip",
	Short: "Retrieve past queries and responses",
	Run:   runClipCmd,
}

func init() {
	rootCmd.AddCommand(clipCmd)
}

func runClipCmd(cmd *cobra.Command, args []string) {
	queryMessage := ""
	clipboardMessage := ""
	err := shared.ExecuteWebSocketCommand("clip", queryMessage, clipboardMessage, "")
	if err != nil {
		fmt.Println(err)
	}
}
