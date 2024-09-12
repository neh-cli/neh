// filename: cmd/o.go

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/coder/websocket"
	"github.com/spf13/cobra"
	"github.com/neh-cli/neh/cmd/shared"
)

var oCmd = &cobra.Command{
    Use:   "o",
    Short: "Send a message to the server",
    Run:   runOCmd,
}

func init() {
    rootCmd.AddCommand(oCmd)
}

func runOCmd(cmd *cobra.Command, args []string) {
    originalMessage := args[0]
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

    shared.HandleWebSocketMessages(ctx, conn, "o", originalMessage, &sync.Map{}, 1, false)
}
