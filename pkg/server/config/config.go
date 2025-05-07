package config

import (
	"time"
)

const (
	DEFAULT_ID                = "chatcmpl-123"
	DEFAULT_MODEL             = "model"
	DEFAULT_FINGERPRINT       = "fp_0123456789"
	CHAT_COMPLETIONS_ENDPOINT = "/v1/chat/completions"
	MODELS_ENDPOINT           = "/v1/models"
	COMPLETIONS_OBJECT        = "chat.completion"
	DEFAULT_PORT              = 8080
	DEFAULT_LISTEN_ADDRESS    = "127.0.0.1"
)

type Config struct {
	Port           int    `json:"port"`
	ChatEndpoint   string `json:"chat_endpoint"`
	ModelsEndpoint string `json:"models_endpoint"`
	Model          string `json:"model"`
	Address        string `json:"address"`
	PerformanceConfig
}

func NewRecommendedConfig() *Config {
	return &Config{
		Port:         DEFAULT_PORT,
		ChatEndpoint: CHAT_COMPLETIONS_ENDPOINT,
		Model:        DEFAULT_MODEL,
		Address:      DEFAULT_LISTEN_ADDRESS,
		PerformanceConfig: PerformanceConfig{
			ITLValue:  1,
			TTFTValue: 2,
		},
	}
}

type PerformanceConfig struct {
	ITLValue  float32 `json:"itl_val"`
	TTFTValue float32 `json:"ttft_val"`
}

func (p *PerformanceConfig) GetITLValue() time.Duration {
	return time.Duration(p.ITLValue * float32(time.Second))
}

func (p *PerformanceConfig) GetTTFTValue() time.Duration {
	return time.Duration(p.TTFTValue * float32(time.Second))
}
