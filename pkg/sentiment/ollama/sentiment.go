package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	pkglog "sb-scanner/pkg/logger"
	"sb-scanner/pkg/sentiment"
)

type OllamaEvaluator struct {
	logger *slog.Logger

	model string
	url   string
	cli   *http.Client
}

func NewOllamaEvaluator(model, url string) *OllamaEvaluator {
	return &OllamaEvaluator{
		logger: pkglog.GetLogger().With("pkg", "OllamaEvaluator"),
		model:  model,
		url:    url,
		cli:    &http.Client{Timeout: 30 * time.Second},
	}
}

// Evaluate calls a local ollama instance to estimate sentiment for the provided text.
// It expects the model to return or include a numeric score between -1.0 and 1.0.
func (e *OllamaEvaluator) Evaluate(ctx context.Context, text string) (sentiment.Sentiment, error) {
	reqBody := chatRequest{
		Model: e.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: text},
		},
		Format: chatFormat{
			Type: "object",
			Properties: map[string]chatFormatProperties{
				"score": {Type: "number"},
			},
			Required: []string{"score"},
		},
		Stream: false,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return sentiment.Sentiment{}, fmt.Errorf("failed to marshal request body: %w", err)
	}
	buf := bytes.NewBuffer(reqBytes)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/chat", e.url), buf)
	if err != nil {
		return sentiment.Sentiment{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.cli.Do(req)
	if err != nil {
		return sentiment.Sentiment{}, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return sentiment.Sentiment{}, fmt.Errorf("failed to read response body: %w", err)
	}
	e.logger.Debug("received response from ollama", "status_code", resp.StatusCode, "response_body", string(respBytes))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return sentiment.Sentiment{}, fmt.Errorf("received non-2xx response from ollama: %d - %s", resp.StatusCode, string(respBytes))
	}
	var chatResp chatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return sentiment.Sentiment{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	var result scoreResult
	if err := json.Unmarshal([]byte(chatResp.Message.Content), &result); err != nil {
		return sentiment.Sentiment{}, fmt.Errorf("failed to unmarshal score from model response: %w", err)
	}

	return sentiment.Sentiment{
		Score: result.Score,
		Model: e.model,
	}, nil
}
