package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"learnforge/internal/domain"
	"net/http"
	"time"
)

type OpenAIClient struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewOpenAIClient(baseURL, apiKey, model string) *OpenAIClient {
	return &OpenAIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *OpenAIClient) ProcessText(ctx context.Context, req *ProcessRequest) (*domain.ProcessResponse, error) {
	prompt := c.buildPrompt(req)
	apiReq := c.createAPIRequest(prompt)
	var resp *domain.ProcessResponse
	var err error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		resp, err = c.makeRequest(ctx, apiReq)
		if err == nil {
			break
		}

		// Check if error is retryable
		if !isRetryableError(err) {
			return nil, err
		}
	}

	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeUpstreamError, "failed to process text with AI", err)
	}

	return resp, nil
}

func (c *OpenAIClient) buildPrompt(req *ProcessRequest) string {
	var promptBuilder bytes.Buffer

	promptBuilder.WriteString("You are an educational content generator. Process the following text and create structured learning content.\n\n")
	promptBuilder.WriteString("Text to process:\n")
	promptBuilder.WriteString(req.Text)
	promptBuilder.WriteString("\n\n")

	switch req.Mode {
	case "flashcards":
		promptBuilder.WriteString("Generate flashcards (question-answer pairs) from this text.\n")
	case "quiz":
		promptBuilder.WriteString("Generate quiz questions with multiple choice answers from this text.\n")
	default:
		promptBuilder.WriteString("Generate a comprehensive lesson with summary, key points, flashcards, and quiz questions.\n")
	}

	if req.Topic != nil {
		promptBuilder.WriteString(fmt.Sprintf("Topic: %s\n", *req.Topic))
	} else {
		promptBuilder.WriteString("Infer the topic from the text and provide your confidence (0.0-1.0).\n")
	}

	if req.Level != nil {
		promptBuilder.WriteString(fmt.Sprintf("Difficulty level: %s\n", *req.Level))
	}

	if req.Language != "" && req.Language != "en" {
		promptBuilder.WriteString(fmt.Sprintf("Language: %s\n", req.Language))
	}

	promptBuilder.WriteString("\n")
	promptBuilder.WriteString("IMPORTANT: Respond ONLY with valid JSON matching this exact schema:\n")
	promptBuilder.WriteString(`{
  "topic": "string",
  "topic_source": "user" or "inferred",
  "topic_confidence": 0.0-1.0,
  "summary": "string",
  "key_points": ["string"],
  "flashcards": [{"q": "string", "a": "string"}],
  "quiz": [{"q": "string", "choices": ["string"], "answer": "string"}]
}`)
	promptBuilder.WriteString("\n\nDo not include any text outside the JSON. Return only the JSON object.")

	return promptBuilder.String()
}

func (c *OpenAIClient) createAPIRequest(prompt string) map[string]interface{} {
	return map[string]interface{}{
		"model": c.model,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"response_format": map[string]string{
			"type": "json_object",
		},
	}
}

func (c *OpenAIClient) makeRequest(ctx context.Context, apiReq map[string]interface{}) (*domain.ProcessResponse, error) {
	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Model string `json:"model"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	var content struct {
		Topic           string      `json:"topic"`
		TopicSource     string      `json:"topic_source"`
		TopicConfidence float64     `json:"topic_confidence"`
		Summary         string      `json:"summary"`
		KeyPoints       []string    `json:"key_points"`
		Flashcards      []domain.Flashcard `json:"flashcards"`
		Quiz            []domain.QuizItem  `json:"quiz"`
	}

	if err := json.Unmarshal([]byte(apiResp.Choices[0].Message.Content), &content); err != nil {
		return nil, fmt.Errorf("failed to parse JSON content: %w", err)
	}

	if content.TopicSource != "user" && content.TopicSource != "inferred" {
		content.TopicSource = "inferred"
	}

	if content.TopicConfidence < 0 {
		content.TopicConfidence = 0
	}
	if content.TopicConfidence > 1 {
		content.TopicConfidence = 1
	}

	response := &domain.ProcessResponse{
		Topic:           content.Topic,
		TopicSource:     content.TopicSource,
		TopicConfidence: content.TopicConfidence,
		Summary:         content.Summary,
		KeyPoints:       content.KeyPoints,
		Flashcards:      content.Flashcards,
		Quiz:            content.Quiz,
		Meta: domain.Meta{
			Model:    apiResp.Model,
			Provider: "openai-compatible",
		},
		CreatedAt: time.Now(),
	}

	return response, nil
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	if err == context.DeadlineExceeded {
		return true
	}

	errStr := err.Error()
	return contains(errStr, "timeout") || contains(errStr, "connection") || contains(errStr, "network")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
			containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

