package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
)

const (
	ollamaURL = "http://ollama:11434"
)

var (
	reqCounter atomic.Uint32
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func askLLMForResponse(ollamaURL, mention string, context []string) (string, error) {
	reqId := reqCounter.Add(1)
	slog.Info("Forwarding request to LLM", "request_id", reqId, "mention", mention, "context", strings.Join(context, ", "))
	contextStr := strings.Join(context, "\n- ")
	prompt := fmt.Sprintf(`You are the owner of a TikTok account that just got mentioned. Please generate a response for the mention.
The Context field may include other comments in the TikTok post to give you more context. Use them if present.

Input Email:
%s

Context:
- %s

Please provide response to the mention without any additional commentary, as it will be sent as-is to the user:`, mention, contextStr)

	reqPayload := OllamaRequest{
		Model:  "phi3:mini",
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}
	slog.Info("Received response from LLM", "request_id", reqId, "response", string(body))

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}
