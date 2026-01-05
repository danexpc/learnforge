package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	webhookURL string
	httpClient *http.Client
}

func NewClient(webhookURL string) *Client {
	return &Client{
		webhookURL: webhookURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Message struct {
	Text        string       `json:"text,omitempty"`
	Blocks      []Block      `json:"blocks,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Block struct {
	Type string `json:"type"`
	Text *Text  `json:"text,omitempty"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Attachment struct {
	Color string `json:"color,omitempty"`
	Text  string `json:"text,omitempty"`
	Title string `json:"title,omitempty"`
}

func (c *Client) SendMessage(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) SendError(ctx context.Context, err error, context map[string]string) error {
	msg := Message{
		Attachments: []Attachment{
			{
				Color: "danger",
				Title: "Error in LearnForge",
				Text:  err.Error(),
			},
		},
	}

	if len(context) > 0 {
		contextText := ""
		for k, v := range context {
			contextText += fmt.Sprintf("*%s*: %s\n", k, v)
		}
		msg.Attachments[0].Text += "\n\n*Context:*\n" + contextText
	}

	return c.SendMessage(ctx, msg)
}

func (c *Client) SendSummary(ctx context.Context, title string, content string) error {
	msg := Message{
		Blocks: []Block{
			{
				Type: "section",
				Text: &Text{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*%s*\n\n%s", title, content),
				},
			},
		},
	}

	return c.SendMessage(ctx, msg)
}

