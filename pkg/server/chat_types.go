package server

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"k8s.io/utils/ptr"
)

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

const (
	OBJECT_TYPE_STREAMING = "chat.completions.chunk"
)

type ChatCompletionsRequestBody struct {
	Messages         []UserMessage          `json:"messages"`
	Model            string                 `json:"model"`
	Store            *string                `json:"store,omitempty"`
	ReasoningEffort  string                 `json:"reasoning_effort"`
	Metadata         map[string]string      `json:"metadata,omitempty"`
	FrequencyPenalty float64                `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]interface{} `json:"logit_bias,omitempty"`

	// Deprecated: Use MaxCompletionTokens instead
	MaxTokens           *int           `json:"max_tokens,omitempty"`
	MaxCompletionTokens *int           `json:"max_completion_tokens,omitempty"`
	N                   *int           `json:"n,omitempty"`
	Modalities          []string       `json:"modalities,omitempty"`
	Prediction          StaticContent  `json:"prediction,omitempty"`
	PresencePenalty     *float64       `json:"presence_penalty,omitempty"`
	Seed                *int           `json:"seed,omitempty"`
	ServiceTier         *string        `json:"service_tier,omitempty"`
	Stop                []string       `json:"stop,omitempty"`
	Stream              bool           `json:"stream,omitempty"`
	StreamOptions       *StreamOptions `json:"stream_options,omitempty"`
	Temperature         *float64       `json:"temperature,omitempty"`
	TopP                *float64       `json:"top_p,omitempty"`
	ParallelToolCalls   bool           `json:"parallel_tool_calls,omitempty"`
	User                string         `json:"user,omitempty"`
}

type StaticContent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type DeveloperMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Name    string `json:"name,omitempty"`
}

type SystemMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Name    string `json:"name,omitempty"`
}

type UserMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Name    string `json:"name,omitempty"`
}

func (m *UserMessage) Validate() bool {
	return m.Content != ""
}

type ChatCompletionsResponse struct {
	ID                string                `json:"id"`
	Choices           []Choice              `json:"choices"`
	Created           int64                 `json:"created"`
	Model             string                `json:"model"`
	ServiceTier       string                `json:"service_tier"`
	SystemFingerprint string                `json:"system_fingerprint"`
	Object            string                `json:"object"`
	Usage             *ChatCompletionsUsage `json:"usage"`
}

type StreamingChatCompletionsResponse struct {
	ID                string                `json:"id"`
	Choices           []StreamingChoice     `json:"choices"`
	Created           int64                 `json:"created"`
	Model             string                `json:"model"`
	ServiceTier       string                `json:"service_tier"`
	SystemFingerprint string                `json:"system_fingerprint"`
	Object            string                `json:"object"`
	Usage             *ChatCompletionsUsage `json:"usage"`
}

type Choice struct {
	FinishReason *string `json:"finish_reason"`
	Index        *int    `json:"index"`
	Message      struct {
		Content *string `json:"content"`
		Refusal *string `json:"refusal"`
		Role    string  `json:"role"`
		// Deprecated: Use ToolCalls instead

	}
}

type Delta struct {
	Content *string `json:"content"`
	Refusal *string `json:"refusal"`
	Role    *string `json:"role"`
}

type StreamingChoice struct {
	FinishReason *string `json:"finish_reason"`
	Index        *int    `json:"index"`
	Delta        Delta   `json:"delta"`
}

func StartedStreamingResponse(
	id string,
	model string,
	fingerprint string,
) StreamingChatCompletionsResponse {
	return StreamingChatCompletionsResponse{
		ID:                id,
		Object:            OBJECT_TYPE_STREAMING,
		Created:           time.Now().Unix(),
		Model:             model,
		SystemFingerprint: fingerprint,
		Usage:             nil,
		Choices: []StreamingChoice{
			{
				FinishReason: nil,
				Index:        ptr.To(0),
				Delta: Delta{
					Content: ptr.To(""),
					Role:    ptr.To("assistant"),
					Refusal: nil,
				},
			},
		},
	}
}

func InProgressStreamingResponse(
	id,
	model string,
	fingerprint string,
	token string,
) StreamingChatCompletionsResponse {
	return StreamingChatCompletionsResponse{
		ID:                id,
		Object:            OBJECT_TYPE_STREAMING,
		Created:           time.Now().Unix(),
		Model:             model,
		SystemFingerprint: fingerprint,
		Usage:             nil,
		Choices: []StreamingChoice{
			{
				FinishReason: nil,
				Index:        ptr.To(0),
				Delta: Delta{
					Content: ptr.To(token),
				},
			},
		},
	}
}

func FinishedStreamingResponse(
	id,
	model string,
	fingerprint string,
	usage *ChatCompletionsUsage,
) StreamingChatCompletionsResponse {

	return StreamingChatCompletionsResponse{
		ID:                id,
		Object:            OBJECT_TYPE_STREAMING,
		Created:           time.Now().Unix(),
		Model:             model,
		SystemFingerprint: fingerprint,
		Usage:             usage,
		Choices: []StreamingChoice{
			{
				FinishReason: ptr.To("stop"),
				Index:        ptr.To(0),
				Delta:        Delta{},
			},
		},
	}
}

func (s *StreamingChatCompletionsResponse) ChunkString() string {
	chunkstring, err := json.Marshal(s)
	if err != nil {
		log.Panic(err)
	}
	return fmt.Sprintf(
		"data: %s\n",
		chunkstring,
	)
}

type ChatCompletionsUsage struct {
	CompletionTokens        int `json:"completion_tokens"`
	PromptTokens            int `json:"prompt_tokens"`
	TotalTokens             int `json:"total_tokens"`
	CompletionTokensDetails struct {
		AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
		AudioTokens              int `json:"audio_tokens"`
		ReasoningTokens          int `json:"reasoning_tokens"`
		RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
	} `json:"completion_tokens_details"`
	PromptTokensDetails struct {
		AudioTokens  int `json:"audio_tokens"`
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
}
