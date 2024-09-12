// filename: cmd/decache.go

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/coder/websocket"
	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
)

var decacheCmd = &cobra.Command{
	Use:   "decache",
	Short: "Remove all query history",
	Long: `This command deletes all previous query history stored on the server.
Use this command to clear any saved interactions with the AI.`,
	Run:   runDecacheCmd,
}

func init() {
	rootCmd.AddCommand(decacheCmd)
}

func runDecacheCmd(cmd *cobra.Command, args []string) {
    commandName := "decache"
    userMessage := ""
    personalAccessToken := os.Getenv("NEH_PERSONAL_ACCESS_TOKEN")
    if personalAccessToken == "" {
        fmt.Println("Please set the environment variable NEH_PERSONAL_ACCESS_TOKEN")
        return
    }

    headers := http.Header{}
    headers.Add("Authorization", fmt.Sprintf("Bearer %s", personalAccessToken))

    ctx := context.Background()
    conn, err := shared.InitializeWebSocketConnection(ctx, personalAccessToken)

    if err != nil {
        fmt.Printf("Failed to connect to WebSocket: %v\n", err)
        return
    }
    defer conn.Close(websocket.StatusInternalError, "Internal error")

    shared.HandleWebSocketMessages(ctx, conn, commandName, userMessage, &sync.Map{}, 1, false)
}
