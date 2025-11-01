package shared

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGetOrCreateDeviceID(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	deviceID1, err := GetOrCreateDeviceID()
	if err != nil {
		t.Fatalf("GetOrCreateDeviceID() error = %v", err)
	}

	if deviceID1 == "" {
		t.Fatal("GetOrCreateDeviceID() returned empty device ID")
	}

	if _, err := uuid.Parse(deviceID1); err != nil {
		t.Fatalf("GetOrCreateDeviceID() returned invalid UUID: %v", err)
	}

	deviceID2, err := GetOrCreateDeviceID()
	if err != nil {
		t.Fatalf("GetOrCreateDeviceID() second call error = %v", err)
	}

	if deviceID1 != deviceID2 {
		t.Errorf("GetOrCreateDeviceID() returned different IDs: %s != %s", deviceID1, deviceID2)
	}
}

func TestGetDeviceIDPath(t *testing.T) {
	path, err := getDeviceIDPath()
	if err != nil {
		t.Fatalf("getDeviceIDPath() error = %v", err)
	}

	if !strings.Contains(path, ".config") {
		t.Errorf("getDeviceIDPath() = %v, want path containing .config", path)
	}

	if !strings.HasSuffix(path, "device_id") {
		t.Errorf("getDeviceIDPath() = %v, want path ending with device_id", path)
	}
}

func TestReadDeviceID(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "device_id")

	t.Run("file does not exist", func(t *testing.T) {
		_, err := readDeviceID(testPath)
		if !os.IsNotExist(err) {
			t.Errorf("readDeviceID() with non-existent file should return os.ErrNotExist, got %v", err)
		}
	})

	t.Run("file exists with valid UUID", func(t *testing.T) {
		expectedID := "550e8400-e29b-41d4-a716-446655440000"
		if err := os.WriteFile(testPath, []byte(expectedID), 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		deviceID, err := readDeviceID(testPath)
		if err != nil {
			t.Fatalf("readDeviceID() error = %v", err)
		}

		if deviceID != expectedID {
			t.Errorf("readDeviceID() = %v, want %v", deviceID, expectedID)
		}
	})

	t.Run("file is empty", func(t *testing.T) {
		emptyPath := filepath.Join(tempDir, "empty_device_id")
		if err := os.WriteFile(emptyPath, []byte(""), 0644); err != nil {
			t.Fatalf("Failed to write empty test file: %v", err)
		}

		_, err := readDeviceID(emptyPath)
		if err == nil {
			t.Error("readDeviceID() with empty file should return error")
		}
	})
}

func TestSaveDeviceID(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "nested", "dir", "device_id")
	testID := "550e8400-e29b-41d4-a716-446655440000"

	if err := saveDeviceID(testPath, testID); err != nil {
		t.Fatalf("saveDeviceID() error = %v", err)
	}

	data, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(data) != testID {
		t.Errorf("Saved device ID = %v, want %v", string(data), testID)
	}
}
