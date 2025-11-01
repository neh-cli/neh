package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func GetOrCreateDeviceID() (string, error) {
	deviceIDPath, err := getDeviceIDPath()
	if err != nil {
		return "", fmt.Errorf("failed to get device ID path: %w", err)
	}

	deviceID, err := readDeviceID(deviceIDPath)
	if err == nil {
		return deviceID, nil
	}

	if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to read device ID: %w", err)
	}

	deviceID = uuid.New().String()
	if err := saveDeviceID(deviceIDPath, deviceID); err != nil {
		return "", fmt.Errorf("failed to save device ID: %w", err)
	}

	return deviceID, nil
}

func getDeviceIDPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "neh")
	return filepath.Join(configDir, "device_id"), nil
}

func readDeviceID(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	deviceID := strings.TrimSpace(string(data))
	if deviceID == "" {
		return "", fmt.Errorf("device ID file is empty")
	}

	return deviceID, nil
}

func saveDeviceID(path, deviceID string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(deviceID), 0644)
}
