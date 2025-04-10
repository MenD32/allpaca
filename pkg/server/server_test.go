package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	allpacaserver "github.com/MenD32/allpaca/pkg/server"
	allpacaconfig "github.com/MenD32/allpaca/pkg/server/config"
)

func TestStreamResponse(t *testing.T) {
	c := allpacaconfig.NewRecommendedConfig()
	s := allpacaserver.NewServer(c)

	body := `{
		"model": "model",
		"messages": [
		{
			"role": "developer",
			"content": "You are a helpful assistant."
		},
		{
			"role": "user",
			"content": "Hello!"
		}
		],
		"stream": true
	}'`
	req := httptest.NewRequest(http.MethodPost, allpacaconfig.CHAT_COMPLETIONS_ENDPOINT, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.HandleChatCompletions(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d, body: %s", w.Code, w.Body.String())
	}
}
