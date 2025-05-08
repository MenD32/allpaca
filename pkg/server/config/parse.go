package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	MIN_TTFT = 0.001
	MIN_ITL  = 0.001
)

func ParseConfigFromFile(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	if config.PerformanceConfig.TTFTValue < MIN_TTFT {
		return nil, fmt.Errorf("TTFT value must be greater than %f", MIN_TTFT)
	}
	if config.PerformanceConfig.ITLValue < MIN_ITL {
		return nil, fmt.Errorf("ITL value must be greater than %f", MIN_ITL)
	}

	if config.ModelsEndpoint == "" {
		config.ModelsEndpoint = MODELS_ENDPOINT
	}
	if config.ChatEndpoint == "" {
		config.ChatEndpoint = CHAT_COMPLETIONS_ENDPOINT
	}
	if config.CompletionsEndpoint == "" {
		config.CompletionsEndpoint = COMPLETIONS_ENDPOINT
	}
	if config.Address == "" {
		config.Address = DEFAULT_LISTEN_ADDRESS
	}
	if config.Port == 0 {
		config.Port = DEFAULT_PORT
	}
	if config.Model == "" {
		config.Model = DEFAULT_MODEL
	}
	if config.PerformanceConfig.ITLValue == 0 {
		config.PerformanceConfig.ITLValue = 0.1
	}
	if config.PerformanceConfig.TTFTValue == 0 {
		config.PerformanceConfig.TTFTValue = 0.2
	}

	return &config, nil
}
