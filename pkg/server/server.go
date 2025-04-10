package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/MenD32/allpaca/pkg/server/config"
	"k8s.io/klog/v2"
)

var (
	DEFAULT_RESPONSE = []string{"Hi", "my", "name", "is", "Joe"}
)

type Server struct {
	config *config.Config
}

func NewServer(c *config.Config) *Server {
	return &Server{
		config: c,
	}
}

func (s *Server) Validate(b *ChatCompletionsRequestBody) bool {
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

func (s *Server) Start() {
	klog.Info("Starting server...")
	http.HandleFunc(s.config.ChatEndpoint, func(w http.ResponseWriter, r *http.Request) {
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

		if !s.Validate(&requestBody) {
			http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
			return
		}

		response := DEFAULT_RESPONSE

		if requestBody.Stream {
			s.streamResponse(w, response, requestBody.StreamOptions)
		} else {
			// TODO: Handle non-streaming response
			// s.WriteResponse(w, response)
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}

	})

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Address, s.config.Port), nil)
	if err != nil {
		klog.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) streamResponse(w http.ResponseWriter, responseTokens []string, options StreamOptions) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	response := make([]StreamingChatCompletionsResponse, len(responseTokens)+2)
	response[0] = StartedStreamingResponse(
		config.DEFAULT_ID,
		s.config.Model,
		config.DEFAULT_FINGERPRINT,
	)

	for i, token := range responseTokens {
		response[i+1] = InProgressStreamingResponse(
			config.DEFAULT_ID,
			s.config.Model,
			config.DEFAULT_FINGERPRINT,
			token,
		)
	}

	var usage *ChatCompletionsUsage = nil
	if options != nil && options.IncludeUsage {
		usage = &ChatCompletionsUsage{
			CompletionTokens: len(responseTokens),
			PromptTokens:     0,
			TotalTokens:      len(responseTokens),
		}
	}

	response[len(response)-1] = FinishedStreamingResponse(
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

	var wg sync.WaitGroup
	wg.Add(len(response) - 2)

	fmt.Fprint(w, response[0].ChunkString())
	flusher.Flush()

	for i, resp := range response[1 : len(response)-1] {
		go func(resp StreamingChatCompletionsResponse, delay time.Duration) {
			defer wg.Done()
			time.Sleep(delay)
			fmt.Fprint(w, resp.ChunkString())
			flusher.Flush()
		}(resp, s.config.TTFTValue+s.config.ITLValue*time.Duration(i))
	}

	wg.Wait()
	fmt.Fprint(w, response[len(response)-1].ChunkString())
	flusher.Flush()

}
