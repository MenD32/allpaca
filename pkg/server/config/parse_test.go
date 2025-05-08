package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/MenD32/allpaca/pkg/server/config"
)

func TestParseConfigFromFile(t *testing.T) {
	t.Run("ValidConfigFile", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "valid_config_*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		configContent := `{
			"port": 8080,
			"chat_endpoint": "/v1/chat/completions",
			"completions_endpoint": "/v1/completions",
			"models_endpoint": "/v1/models",
			"model": "model",
			"address": "127.0.0.1",
			"itl_val": 1,
			"ttft_val": 2
		}`
		if _, err := tempFile.WriteString(configContent); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tempFile.Close()

		c, err := config.ParseConfigFromFile(tempFile.Name())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if c.Port != 8080 ||
			c.ChatEndpoint != "/v1/chat/completions" ||
			c.Model != "model" ||
			c.Address != "127.0.0.1" ||
			c.PerformanceConfig.GetITLValue() != time.Second ||
			c.PerformanceConfig.GetTTFTValue() != time.Second*2 {
			t.Errorf("Unexpected config values: %+v", c)
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := config.ParseConfigFromFile("non_existent_file.json")
		if err == nil {
			t.Fatal("Expected an error for non-existent file, got nil")
		}
	})

	t.Run("InvalidJSONFile", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "invalid_config_*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		invalidContent := `{"field1": "value1", "field2":`
		if _, err := tempFile.WriteString(invalidContent); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tempFile.Close()

		_, err = config.ParseConfigFromFile(tempFile.Name())
		if err == nil {
			t.Fatal("Expected an error for invalid JSON, got nil")
		}
	})
}
