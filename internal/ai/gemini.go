package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"learnforge/internal/domain"
)

type GeminiClient struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewGeminiClient(apiKey, model string) *GeminiClient {
	if model == "" {
		model = "gemini-2.0-flash-exp"
	}
	return &GeminiClient{
		baseURL: "https://generativelanguage.googleapis.com",
		apiKey:  apiKey,
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *GeminiClient) ProcessText(ctx context.Context, req *ProcessRequest) (*domain.ProcessResponse, error) {
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

		if !isRetryableError(err) {
			return nil, err
		}
	}

	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeUpstreamError, "failed to process text with Gemini", err)
	}

	return resp, nil
}

func (c *GeminiClient) buildPrompt(req *ProcessRequest) string {
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

func (c *GeminiClient) createAPIRequest(prompt string) map[string]interface{} {
	return map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":      0.7,
			"responseMimeType": "application/json",
		},
	}
}

func (c *GeminiClient) makeRequest(ctx context.Context, apiReq map[string]interface{}) (*domain.ProcessResponse, error) {
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)

	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		Model string `json:"model,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	var content struct {
		Topic           string             `json:"topic"`
		TopicSource     string             `json:"topic_source"`
		TopicConfidence float64            `json:"topic_confidence"`
		Summary         string             `json:"summary"`
		KeyPoints       []string           `json:"key_points"`
		Flashcards      []domain.Flashcard `json:"flashcards"`
		Quiz            []domain.QuizItem  `json:"quiz"`
	}

	responseText := apiResp.Candidates[0].Content.Parts[0].Text
	if err := json.Unmarshal([]byte(responseText), &content); err != nil {
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

	modelName := c.model
	if apiResp.Model != "" {
		modelName = apiResp.Model
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
			Model:    modelName,
			Provider: "gemini",
		},
		CreatedAt: time.Now(),
	}

	return response, nil
}
