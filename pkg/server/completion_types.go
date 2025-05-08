package server

import (
	"encoding/json"
	"fmt"
	"log"
)

// CompletionRequest defines the parameters for a completion API request
type CompletionRequest struct {
	Model            string             `json:"model"`
	Prompt           any                `json:"prompt,omitempty"`
	Suffix           string             `json:"suffix,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	N                int                `json:"n,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	LogProbs         int                `json:"logprobs,omitempty"`
	Echo             bool               `json:"echo,omitempty"`
	Stop             any                `json:"stop,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	BestOf           int                `json:"best_of,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
}

func (c *CompletionRequest) GetPrompts() []string {
	switch v := c.Prompt.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case []any:
		prompts := make([]string, len(v))
		for i, prompt := range v {
			if str, ok := prompt.(string); ok {
				prompts[i] = str
			}
		}
		return prompts
	default:
		return []string{}
	}
}

// CompletionChoice represents a completion choice returned by the API
type CompletionChoice struct {
	Text         string    `json:"text"`
	Index        int       `json:"index"`
	LogProbs     *LogProbs `json:"logprobs"`
	FinishReason string    `json:"finish_reason"`
}

// LogProbs contains the log probabilities for tokens
type LogProbs struct {
	Tokens        []string             `json:"tokens"`
	TokenLogProbs []float64            `json:"token_logprobs"`
	TopLogProbs   []map[string]float64 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

// CompletionResponse represents a response from the completion API
type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// CompletionStreamResponse represents a streaming response chunk
type CompletionStreamResponse struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int64                    `json:"created"`
	Model   string                   `json:"model"`
	Choices []CompletionStreamChoice `json:"choices"`
}

// CompletionStreamChoice represents a single choice in a streaming response
type CompletionStreamChoice struct {
	Text         string    `json:"text"`
	Index        int       `json:"index"`
	LogProbs     *LogProbs `json:"logprobs"`
	FinishReason string    `json:"finish_reason"`
}

func (s *CompletionStreamResponse) ChunkString() string {
	chunkstring, err := json.Marshal(s)
	if err != nil {
		log.Panic(err)
	}
	return fmt.Sprintf(
		"data: %s\n\n",
		chunkstring,
	)
}
