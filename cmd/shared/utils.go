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

func ExecuteWebSocketCommand(command, message string) error {
	personalAccessToken, err := getPersonalAccessToken()
	if err != nil {
		return err
	}

	headers := createAuthorizationHeader(personalAccessToken)

	ctx := context.Background()
	conn, err := initializeWebSocketConnection(ctx, headers)
	if err != nil {
		return fmt.Errorf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close(websocket.StatusInternalError, "Internal error")

	handleWebSocketMessages(ctx, conn, command, message, &sync.Map{})
	return nil
}

func initializeWebSocketConnection(ctx context.Context, headers http.Header) (*websocket.Conn, error) {
	wsURL := getWSUrl()

	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: headers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to dial websocket: %w", err)
	}

	uuidStr := uuid.New().String()
	if err := subscribeToChannel(ctx, conn, uuidStr); err != nil {
		conn.Close(websocket.StatusInternalError, "Subscription failed")
		return nil, fmt.Errorf("failed to subscribe to channel: %w", err)
	}

	return conn, nil
}

func subscribeToChannel(ctx context.Context, conn *websocket.Conn, uuid string) error {
	identifier, err := createIdentifier(uuid)
	if err != nil {
		return fmt.Errorf("failed to marshal identifier: %w", err)
	}

	content := createSubscriptionContent(identifier)

	if err := wsjson.Write(ctx, conn, content); err != nil {
		return fmt.Errorf("failed to write subscription message: %w", err)
	}
	return nil
}

func createIdentifier(uuid string) (string, error) {
	identifier := map[string]interface{}{
		"channel": "LargeLanguageModelQueryChannel",
		"uuid":    uuid,
	}
	identifierJSON, err := json.Marshal(identifier)
	if err != nil {
		return "", err
	}
	return string(identifierJSON), nil
}

func createSubscriptionContent(identifier string) map[string]interface{} {
	return map[string]interface{}{
		"command":    "subscribe",
		"identifier": identifier,
	}
}

func getPersonalAccessToken() (string, error) {
	personalAccessToken := os.Getenv("NEH_PERSONAL_ACCESS_TOKEN")
	if personalAccessToken == "" {
		return "", fmt.Errorf("Please set the environment variable NEH_PERSONAL_ACCESS_TOKEN")
	}
	return personalAccessToken, nil
}

func createAuthorizationHeader(token string) http.Header {
	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return headers
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

func HandleActionCableMessages(conn *websocket.Conn, command string, message map[string]interface{}, originalMessage string) {
	switch message["type"] {
	case "welcome":
		handleWelcomeMessage(conn)
	case "confirm_subscription":
		handleConfirmSubscriptionMessage(conn, message, command, originalMessage)
	case "ping":
		// do nothing
	case "disconnect":
		handleDisconnectMessage(conn, message)
	default:
		handleUnknownMessageType(conn, message)
	}
}

func handleWelcomeMessage(conn *websocket.Conn) {
	uuid := uuid.New().String()
	subscribe(conn, uuid)
}

func handleConfirmSubscriptionMessage(conn *websocket.Conn, message map[string]interface{}, command, originalMessage string) {
	identifier, ok := message["identifier"].(string)
	if !ok {
		fmt.Println("Error: 'identifier' field is missing or not a string")
		return
	}

	onSubscribed(identifier, command, originalMessage)
}

func handleDisconnectMessage(conn *websocket.Conn, message map[string]interface{}) {
	fmt.Printf("Connection has been disconnected. Reason: %s\n", message["reason"])
	conn.Close(websocket.StatusNormalClosure, "Normal closure")
}

func handleUnknownMessageType(conn *websocket.Conn, message map[string]interface{}) {
	fmt.Printf("unknown message type in handleActionCableMessages: %s. Closing connection.\n", message["type"])
	conn.Close(websocket.StatusNormalClosure, "Normal closure")
}

func handleWebSocketMessages(ctx context.Context, conn *websocket.Conn, command string, originalMessage string, messagePool *sync.Map) {
	var expectedSequenceNumber uint = 1

	for {
		var message map[string]interface{}
		err := wsjson.Read(ctx, conn, &message)

		if err != nil {
			fmt.Println("")
			break
		} else if message["type"] != nil {
			HandleActionCableMessages(conn, command, message, originalMessage)
		} else {
			HandleBroadcastedMessages(conn, message, messagePool, &expectedSequenceNumber)
		}
	}
}

func onSubscribed(identifier string, command string, message string) {
	personalAccessToken := os.Getenv("NEH_PERSONAL_ACCESS_TOKEN")
	if os.Getenv("NEH_DEBUG") == "t" {
		fmt.Printf("Personal Access Token: %s\n", personalAccessToken)
	}

	httpURL := getNehServerEndpoint(command)

	identifierMap, err := unmarshalIdentifier(identifier)
	if err != nil {
		fmt.Printf("Failed to unmarshal identifier: %v\n", err)
		return
	}

	uuid, ok := identifierMap["uuid"].(string)
	if !ok {
		fmt.Println("Error: 'uuid' field is missing or not a string")
		return
	}

	reqBody, err := createRequestBody(message, uuid, personalAccessToken)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	if err := sendHttpRequest(httpURL, reqBody, personalAccessToken); err != nil {
		fmt.Printf("Failed to send HTTP request: %v\n", err)
	}
}

func unmarshalIdentifier(identifier string) (map[string]interface{}, error) {
	var identifierMap map[string]interface{}
	if err := json.Unmarshal([]byte(identifier), &identifierMap); err != nil {
		return nil, err
	}
	return identifierMap, nil
}

func createRequestBody(message, uuid, token string) ([]byte, error) {
	reqBody := map[string]interface{}{
		"message": message,
		"uuid":    uuid,
		"token":   token,
	}
	return json.Marshal(reqBody)
}

func sendHttpRequest(url string, body []byte, token string) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handleHttpResponse(resp)
}

func handleHttpResponse(resp *http.Response) error {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(respBody, &responseBody); err != nil {
		return err
	}

	if responseBody["message"] != nil {
		fmt.Printf("%s", responseBody["message"])
	}
	return nil
}

func getNehServerEndpoint(command string) string {
	if os.Getenv("WORKING_ON_LOCALHOST") == "t" {
		developmentEndpoint := os.Getenv("NEH_SERVER_ENDPOINT_DEVELOPMENT")
		if developmentEndpoint == "" {
			panic("The environment variable NEH_SERVER_ENDPOINT_DEVELOPMENT is not set")
		}
		return fmt.Sprintf("%s%s", developmentEndpoint, command)
	}
	return fmt.Sprintf("https://yoryo-app.onrender.com/api/neh/%s", command)
}

func HandleBroadcastedMessages(conn *websocket.Conn, message map[string]interface{}, messagePool *sync.Map, expectedSequenceNumber *uint) {
	messageType, message, err := extractMessageType(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch messageType {
	case "output":
		handleOutputMessage(message, messagePool, expectedSequenceNumber)
	case "error":
		fmt.Printf("Error message received: %v\n", message["body"])
	case "worker_done":
		conn.Close(websocket.StatusNormalClosure, "Normal closure")
	default:
		fmt.Printf("Unknown message type in handleBroadcastedMessages: %v\n", messageType)
	}
}

func extractMessageType(message map[string]interface{}) (string, map[string]interface{}, error) {
	if messageType, ok := message["type"].(string); ok {
		return messageType, message, nil
	}

	innerMessage, ok := message["message"].(map[string]interface{})
	if !ok {
		return "", nil, fmt.Errorf("Error: 'type' field is missing or not a string")
	}

	messageType, ok := innerMessage["type"].(string)
	if !ok {
		return "", nil, fmt.Errorf("Error: 'type' field is missing or not a string in inner message")
	}

	return messageType, innerMessage, nil
}

func handleOutputMessage(message map[string]interface{}, messagePool *sync.Map, expectedSequenceNumber *uint) {
	if sequenceNumber, ok := message["sequence_number"].(float64); ok {
		messagePool.Store(uint(sequenceNumber), message["body"].(string))
		processMessageInOrder(messagePool, expectedSequenceNumber)
	} else {
		fmt.Println("Error: 'sequence_number' field is missing or not a float64")
	}
}

func processMessageInOrder(messagePool *sync.Map, expectedSequenceNumber *uint) {
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
