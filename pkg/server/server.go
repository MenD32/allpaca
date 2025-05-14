package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MenD32/allpaca/pkg/server/config"
	"k8s.io/klog/v2"
)

type Server struct {
	config *config.Config
}

func NewServer(c *config.Config) *Server {
	return &Server{
		config: c,
	}
}

func (s *Server) ValidateChatCompletionsRequestBody(b *ChatCompletionsRequestBody) bool {
	if b.Model != s.config.Model {
		klog.Infof("Model %s is not supported", b.Model)
		return false
	}
	if len(b.Messages) == 0 {
		klog.Info("Messages cannot be empty")
		return false
	}
	for _, m := range b.Messages {
		if !m.Validate() {
			return false
		}
	}
	return true
}

func (s *Server) ValidateCompletionRequest(b *CompletionRequest) bool {
	if b.Model == "" {
		klog.Info("Model cannot be empty")
		return false
	}

	if b.Model != s.config.Model {
		klog.Infof("Model %s is not supported", b.Model)
		return false
	}

	if b.Prompt == nil {
		klog.Info("Prompt cannot be empty")
		return false
	}

	return true
}

func (s *Server) Start() {
	klog.Infof("Server configuration: %+v", s.config)
	klog.Infof("TTFT: %d", s.config.GetTTFTValue())
	klog.Infof("ITL: %d", s.config.GetITLValue())
	klog.Infof("Starting server on %s:%d...", s.config.Address, s.config.Port)
	klog.Infof("Chat endpoint: %s", s.config.ChatEndpoint)
	http.HandleFunc(s.config.ChatEndpoint, s.HandleChatCompletions)
	klog.Infof("Completions endpoint: %s", s.config.CompletionsEndpoint)
	http.HandleFunc(s.config.CompletionsEndpoint, s.HandleCompletions)
	klog.Infof("Models endpoint: %s", s.config.ModelsEndpoint)
	http.HandleFunc(s.config.ModelsEndpoint, s.HandleModels)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		klog.Infof("Request path: %s", r.URL.Path)
		http.Error(w, "Not found", http.StatusNotFound)
	})

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Address, s.config.Port), nil)
	if err != nil {
		klog.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) HandleModels(w http.ResponseWriter, r *http.Request) {
	klog.Info("Received request for models")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	models := []Model{
		{
			ID:      s.config.Model,
			Object:  s.config.Model,
			Created: time.Now().Unix(),
			OwnedBy: "allpaca",
		},
	}
	response := map[string]interface{}{
		"object": "list",
		"data":   models,
	}

	responsebytes, err := json.Marshal(response)
	if err != nil {
		klog.Errorf("Failed to marshal models response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	klog.Infof("Models response: %s", response)
}

func (s *Server) HandleChatCompletions(w http.ResponseWriter, r *http.Request) {
	klog.Info("Received request for chat completions")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestBody ChatCompletionsRequestBody

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		klog.Errorf("Failed to marshal request body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	klog.Infof("Request Body: %s", requestBodyJSON)

	if !s.ValidateChatCompletionsRequestBody(&requestBody) {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		return
	}

	var response []string
	var prompt_token_count, completion_token_count, total_token_count int
	var usage *ChatCompletionsUsage = nil

	for _, message := range requestBody.Messages {
		response = append(response, strings.Split(message.Content, " ")...)
		prompt_token_count += len(strings.Split(message.Content, " "))
	}
	completion_token_count = len(response)
	total_token_count = prompt_token_count + completion_token_count

	usage = &ChatCompletionsUsage{
		CompletionTokens: completion_token_count,
		PromptTokens:     prompt_token_count,
		TotalTokens:      total_token_count,
	}

	if requestBody.Stream {
		s.streamResponse(w, response, requestBody.StreamOptions, usage)
	} else {
		// TODO: Handle non-streaming response
		// s.WriteResponse(w, response)
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}

}

func (s *Server) HandleCompletions(w http.ResponseWriter, r *http.Request) {
	klog.Info("Received request for chat completions")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestBody CompletionRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		klog.Errorf("Failed to marshal request body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	klog.Infof("Request Body: %s", requestBodyJSON)

	if !s.ValidateCompletionRequest(&requestBody) {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		return
	}

	var response []string

	for _, message := range requestBody.GetPrompts() {
		response = append(response, strings.Split(message, " ")...)
	}

	if requestBody.Stream {
		s.streamResponseLegacy(w, response)
	} else {
		// TODO: Handle non-streaming response
		// s.WriteResponse(w, response)
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}

}

func (s *Server) streamResponseLegacy(w http.ResponseWriter, responseTokens []string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	response := make([]CompletionStreamResponse, len(responseTokens))

	for i, token := range responseTokens {
		finishReason := ""
		if i == len(responseTokens)-1 {
			finishReason = "stop"
		}
		response[i] = CompletionStreamResponse{
			ID:      config.DEFAULT_ID,
			Object:  "text_completion",
			Created: time.Now().Unix(),
			Model:   s.config.Model,
			Choices: []CompletionStreamChoice{
				{
					Text:         token+" ",
					Index:        i,
					LogProbs:     nil,
					FinishReason: finishReason,
				},
			},
		}
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	start_time := time.Now().Add(s.config.GetTTFTValue())

	time.Sleep(time.Until(start_time))

	for i, resp := range response {
		time.Sleep(time.Until(start_time.Add(s.config.GetITLValue() * time.Duration(i))))
		fmt.Fprint(w, resp.ChunkString())
		flusher.Flush()
	}

	fmt.Fprint(w, "data: [DONE]\n\n")
	flusher.Flush()

	w.Header().Set("HTTP/2", "200")
}

func (s *Server) streamResponse(w http.ResponseWriter, responseTokens []string, options *StreamOptions, usage *ChatCompletionsUsage) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	response := make([]StreamingChatCompletionsResponse, len(responseTokens))
	start_response := StartedStreamingResponse(
		config.DEFAULT_ID,
		s.config.Model,
		config.DEFAULT_FINGERPRINT,
	)

	for i, token := range responseTokens {
		response[i] = InProgressStreamingResponse(
			config.DEFAULT_ID,
			s.config.Model,
			config.DEFAULT_FINGERPRINT,
			token+" ",
		)
	}

	if options != nil && options.IncludeUsage {
		usage = nil
	}
	finish_response := FinishedStreamingResponse(
		config.DEFAULT_ID,
		s.config.Model,
		config.DEFAULT_FINGERPRINT,
		usage,
	)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	start_time := time.Now().Add(s.config.GetTTFTValue())

	fmt.Fprint(w, start_response.ChunkString())
	flusher.Flush()

	time.Sleep(time.Until(start_time))

	for i, resp := range response {
		time.Sleep(time.Until(start_time.Add(s.config.GetITLValue() * time.Duration(i))))
		fmt.Fprint(w, resp.ChunkString())
		flusher.Flush()
	}

	fmt.Fprint(w, finish_response.ChunkString())
	flusher.Flush()

	fmt.Fprint(w, "data: [DONE]\n\n")
	flusher.Flush()

	w.Header().Set("HTTP/2", "200")
}
