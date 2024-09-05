package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var OCmd = &cobra.Command{
    Use:   "o",
    Short: "Send a message to the server",
    Run:   runOCmd,
}

func init() {
    rootCmd.AddCommand(OCmd)
}

func runOCmd(cmd *cobra.Command, args []string) {
    originalMessage := args[0]
    personalAccessToken := os.Getenv("NEH_PERSONAL_ACCESS_TOKEN")
    if personalAccessToken == "" {
        fmt.Println("Please set the environment variable NEH_PERSONAL_ACCESS_TOKEN.")
        return
    }

    wsURL := getWSUrl()
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
    subscribe(conn, uuid)

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
            handleActionCableMessages(conn, message, originalMessage, &requestSent)
        } else {
            handleBroadcastedMessages(conn, message, &messagePool, &expectedSequenceNumber)
        }
    }
}

func getWSUrl() string {
    if os.Getenv("WORKING_ON_LOCALHOST") != "" {
        return "ws://localhost:6060/cable"
    }
    return "wss://yoryo-app.onrender.com/cable"
}

func subscribe(conn *websocket.Conn, uuid string) {
    identifier := map[string]interface{}{
        "channel": "LargeLanguageModelQueryChannel",
        "uuid":    uuid,
    }
    identifierJSON, _ := json.Marshal(identifier)
    content := map[string]interface{}{
        "command":    "subscribe",
        "identifier": string(identifierJSON),
    }
    wsjson.Write(context.Background(), conn, content)
}

func handleActionCableMessages(conn *websocket.Conn, message map[string]interface{}, originalMessage string, requestSent *bool) {
    switch message["type"] {
    case "welcome":
        uuid := uuid.New().String()
        subscribe(conn, uuid)
    case "confirm_subscription":
        if !*requestSent {
            if identifier, ok := message["identifier"].(string); ok {
                onSubscribed(identifier, originalMessage)
                *requestSent = true
            } else {
                fmt.Println("Error: 'identifier' field is missing or not a string")
            }
        }
    case "ping":
        // do nothing
    case "disconnect":
        fmt.Printf("Connection has been disconnected. Reason: %s\n", message["reason"])
        conn.Close(websocket.StatusNormalClosure, "Normal closure")
    default:
        fmt.Printf("unknown message type in handleActionCableMessages: %s. Closing connection.\n", message["type"])
        conn.Close(websocket.StatusNormalClosure, "Normal closure")
    }
}

func handleBroadcastedMessages(conn *websocket.Conn, message map[string]interface{}, messagePool *sync.Map, expectedSequenceNumber *int) {
    messageType, ok := message["type"].(string)
    if !ok {
        if innerMessage, ok := message["message"].(map[string]interface{}); ok {
            messageType, ok = innerMessage["type"].(string)
            if !ok {
                fmt.Println("Error: 'type' field is missing or not a string in inner message")
                return
            }
            message = innerMessage
        } else {
            fmt.Println("Error: 'type' field is missing or not a string")
            return
        }
    }

    switch messageType {
    case "output":
        if sequenceNumber, ok := message["sequence_number"].(float64); ok {
            messagePool.Store(int(sequenceNumber), message["body"].(string))
            processMessageInOrder(messagePool, expectedSequenceNumber)
        } else {
            fmt.Println("Error: 'sequence_number' field is missing or not a float64")
        }
    case "error":
        fmt.Printf("Error message received: %v\n", message["body"])
    case "worker_done":
        conn.Close(websocket.StatusNormalClosure, "Normal closure")
    default:
        fmt.Printf("Unknown message type in handleBroadcastedMessages: %v\n", messageType)
    }
}

func processMessageInOrder(messagePool *sync.Map, expectedSequenceNumber *int) {
    for {
        if value, ok := messagePool.Load(*expectedSequenceNumber); ok {
            fmt.Print(value)
            messagePool.Delete(*expectedSequenceNumber)
            *expectedSequenceNumber++
        } else {
            break
        }
    }
}

func onSubscribed(identifier string, message string) {
    personalAccessToken := os.Getenv("NEH_PERSONAL_ACCESS_TOKEN")
    httpURL := getHttpUrl()

    var identifierMap map[string]interface{}
    if err := json.Unmarshal([]byte(identifier), &identifierMap); err != nil {
        fmt.Printf("Failed to unmarshal identifier: %v\n", err)
        return
    }

    reqBody := map[string]interface{}{
        "message": message,
        "uuid":    identifierMap["uuid"],
        "token":   personalAccessToken,
    }
    body, err := json.Marshal(reqBody)
    if err != nil {
        fmt.Printf("Failed to marshal request body: %v\n", err)
        return
    }

    req, err := http.NewRequest("POST", httpURL, bytes.NewBuffer(body))
    if err != nil {
        fmt.Printf("Failed to create HTTP request: %v\n", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", personalAccessToken))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("Failed to send HTTP request: %v\n", err)
        return
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)
    var responseBody map[string]interface{}
    json.Unmarshal(respBody, &responseBody)

    if responseBody["message"] != nil {
        fmt.Println(responseBody["message"])
    }
}

func getHttpUrl() string {
    if os.Getenv("WORKING_ON_LOCALHOST") != "" {
        return "http://localhost:6060/api/neh/o"
    }
    return "https://yoryo-app.onrender.com/api/neh/o"
}
