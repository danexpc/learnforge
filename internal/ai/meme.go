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

// ImgflipMemeGenerator is a free meme generator using Imgflip API
type ImgflipMemeGenerator struct {
	httpClient *http.Client
}

// NewImgflipMemeGenerator creates a new Imgflip meme generator
func NewImgflipMemeGenerator() *ImgflipMemeGenerator {
	return &ImgflipMemeGenerator{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateMeme generates a meme using Imgflip's free API
// Popular meme templates: 181913649 (Drake), 87743020 (Two Buttons), 112126428 (Distracted Boyfriend)
func (g *ImgflipMemeGenerator) GenerateMeme(ctx context.Context, topic, question string) (string, error) {
	// Select a random educational meme template
	templates := []struct {
		ID   int
		Name string
	}{
		{181913649, "Drake"},           // Good for comparisons
		{87743020, "Two Buttons"},      // Good for choices
		{112126428, "Distracted Boyfriend"}, // Good for focus
		{129242436, "Change My Mind"},  // Good for opinions
		{438680, "Batman Slapping"},    // Good for reactions
	}

	// Use a simple template based on topic hash for consistency
	templateID := templates[0].ID // Default to Drake template

	// Create meme text based on topic and question
	var topText, bottomText string
	if question != "" {
		// Use question as top text, topic as bottom
		topText = truncate(question, 100)
		bottomText = truncate(fmt.Sprintf("Learning about %s", topic), 100)
	} else {
		topText = truncate(fmt.Sprintf("When you understand %s", topic), 100)
		bottomText = "But you still need to study"
	}

	// Call Imgflip API
	apiURL := "https://api.imgflip.com/caption_image"
	data := url.Values{}
	data.Set("template_id", fmt.Sprintf("%d", templateID))
	data.Set("username", "imgflip_hubot") // Public demo account
	data.Set("password", "imgflip_hubot")  // Public demo account
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

// truncate truncates a string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

