package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ImgflipMemeGenerator struct {
	httpClient *http.Client
}

func NewImgflipMemeGenerator() *ImgflipMemeGenerator {
	return &ImgflipMemeGenerator{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (g *ImgflipMemeGenerator) GenerateMeme(ctx context.Context, topic, question string) (string, error) {
	templates := []struct {
		ID   int
		Name string
	}{
		{181913649, "Drake"},
		{87743020, "Two Buttons"},
		{112126428, "Distracted Boyfriend"},
		{129242436, "Change My Mind"},
		{438680, "Batman Slapping"},
	}

	templateID := templates[0].ID

	var topText, bottomText string
	if question != "" {
		topText = truncate(question, 100)
		bottomText = truncate(fmt.Sprintf("Learning about %s", topic), 100)
	} else {
		topText = truncate(fmt.Sprintf("When you understand %s", topic), 100)
		bottomText = "But you still need to study"
	}

	apiURL := "https://api.imgflip.com/caption_image"
	data := url.Values{}
	data.Set("template_id", fmt.Sprintf("%d", templateID))
	data.Set("username", "imgflip_hubot")
	data.Set("password", "imgflip_hubot")
	data.Set("text0", topText)
	data.Set("text1", bottomText)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Success bool `json:"success"`
		Data    struct {
			URL string `json:"url"`
		} `json:"data"`
		ErrorMessage string `json:"error_message,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return "", fmt.Errorf("imgflip API error: %s", apiResp.ErrorMessage)
	}

	if apiResp.Data.URL == "" {
		return "", fmt.Errorf("no image URL in response")
	}

	return apiResp.Data.URL, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

