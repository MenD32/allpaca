package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	MIN_TTFT = 0.01
	MIN_ITL  = 0.01
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

	return &config, nil
}
