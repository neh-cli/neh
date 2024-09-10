// filename: cmd/decache.go

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
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
    personalAccessToken := os.Getenv("NEH_PERSONAL_ACCESS_TOKEN")
    if personalAccessToken == "" {
        fmt.Println("Please set the environment variable NEH_PERSONAL_ACCESS_TOKEN.")
        return
    }

    wsURL := shared.GetWSUrl()
    headers := http.Header{}
    headers.Add("Authorization", fmt.Sprintf("Bearer %s", personalAccessToken))

    conn, _, err := websocket.Dial(context.Background(), wsURL, &websocket.DialOptions{
        HTTPHeader: headers,
    })
    if err != nil {
        fmt.Printf("Failed to connect to WebSocket: %v\n", err)
        return
    }
    defer conn.Close(websocket.StatusInternalError, "Internal error")

    uuid := uuid.New().String()
    shared.Subscribe(conn, uuid)

    var messagePool sync.Map
    expectedSequenceNumber := 1
    requestSent := false

    for {
        var message map[string]interface{}
        err := wsjson.Read(context.Background(), conn, &message)

        if err != nil {
            fmt.Println("")
            break
        } else if message["type"] != nil {
            shared.HandleActionCableMessages(conn, "decache", message, "", &requestSent)
        } else {
            shared.HandleBroadcastedMessages(conn, message, &messagePool, &expectedSequenceNumber)
        }
    }
}