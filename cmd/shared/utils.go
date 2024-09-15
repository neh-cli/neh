// filename: cmd/shared/utils.go

package shared

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
)

func ExecuteWebSocketCommand(command, message string, waitForResponse bool) error {
    personalAccessToken := os.Getenv("NEH_PERSONAL_ACCESS_TOKEN")
    if personalAccessToken == "" {
        return fmt.Errorf("Please set the environment variable NEH_PERSONAL_ACCESS_TOKEN")
    }

    headers := http.Header{}
    headers.Add("Authorization", fmt.Sprintf("Bearer %s", personalAccessToken))

    ctx := context.Background()
    conn, err := InitializeWebSocketConnection(ctx, personalAccessToken)
    if err != nil {
        return fmt.Errorf("Failed to connect to WebSocket: %v", err)
    }
    defer conn.Close(websocket.StatusInternalError, "Internal error")

    HandleWebSocketMessages(ctx, conn, command, message, &sync.Map{}, waitForResponse)
    return nil
}

func GetWSUrl() string {
    if os.Getenv("WORKING_ON_LOCALHOST") != "" {
        return "ws://localhost:6060/cable"
    }
    return "wss://yoryo-app.onrender.com/cable"
}

func InitializeWebSocketConnection(ctx context.Context, personalAccessToken string) (*websocket.Conn, error) {
    wsURL := GetWSUrl()
    headers := http.Header{}
    headers.Add("Authorization", fmt.Sprintf("Bearer %s", personalAccessToken))

    conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
        HTTPHeader: headers,
    })
    if err != nil {
        return nil, err
    }

    uuidStr := uuid.New().String()
    Subscribe(conn, uuidStr)
    return conn, nil
}

func Subscribe(conn *websocket.Conn, uuid string) {
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

func HandleActionCableMessages(conn *websocket.Conn, command string, message map[string]interface{}, originalMessage string, requestSent *bool) {
    switch message["type"] {
    case "welcome":
        uuid := uuid.New().String()
        Subscribe(conn, uuid)
    case "confirm_subscription":
        if !*requestSent {
            if identifier, ok := message["identifier"].(string); ok {
                onSubscribed(identifier, command, originalMessage)
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

func HandleWebSocketMessages(ctx context.Context, conn *websocket.Conn, command string, originalMessage string, messagePool *sync.Map, requestSent bool) {
    expectedSequenceNumber := 1

    for {
        var message map[string]interface{}
        err := wsjson.Read(ctx, conn, &message)

        if err != nil {
            fmt.Println("")
            break
        } else if message["type"] != nil {
            HandleActionCableMessages(conn, command, message, originalMessage, &requestSent)
        } else {
            HandleBroadcastedMessages(conn, message, messagePool, &expectedSequenceNumber)
        }
    }
}

func onSubscribed(identifier string, command string, message string) {
    personalAccessToken := os.Getenv("NEH_PERSONAL_ACCESS_TOKEN")
    httpURL := getHttpUrl(command)

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
        fmt.Printf("%s", responseBody["message"])
    }
}

func getHttpUrl(command string) string {
    if os.Getenv("WORKING_ON_LOCALHOST") != "" {
        return fmt.Sprintf("http://localhost:6060/api/neh/%s", command)
    }
    return fmt.Sprintf("https://yoryo-app.onrender.com/api/neh/%s", command)
}

func HandleBroadcastedMessages(conn *websocket.Conn, message map[string]interface{}, messagePool *sync.Map, expectedSequenceNumber *int) {
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
