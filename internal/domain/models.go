package domain

import "time"

// ProcessRequest represents the incoming request to process text
type ProcessRequest struct {
	Text          string  `json:"text"`
	Mode          string  `json:"mode,omitempty"` // lesson, flashcards, quiz
	Topic         *string `json:"topic,omitempty"`
	Level         *string `json:"level,omitempty"` // beginner, intermediate, advanced
	Language      string  `json:"language,omitempty"`
	IdempotencyKey *string `json:"idempotency_key,omitempty"`
}

// ProcessResponse represents the structured learning content response
type ProcessResponse struct {
	ID              string    `json:"id"`
	Topic           string    `json:"topic"`
	TopicSource     string    `json:"topic_source"`     // user, inferred
	TopicConfidence float64   `json:"topic_confidence"` // 0.0-1.0
	Summary         string    `json:"summary"`
	KeyPoints       []string  `json:"key_points"`
	Flashcards      []Flashcard `json:"flashcards"`
	Quiz            []QuizItem  `json:"quiz"`
	Meta            Meta      `json:"meta"`
	CreatedAt       time.Time `json:"created_at"`
}

// Flashcard represents a question-answer pair
type Flashcard struct {
	Q string `json:"q"`
	A string `json:"a"`
}

// QuizItem represents a quiz question with multiple choice
type QuizItem struct {
	Q       string   `json:"q"`
	Choices []string `json:"choices"`
	Answer  string   `json:"answer"`
}

// Meta contains processing metadata
type Meta struct {
	Model        string  `json:"model"`
	Provider     string  `json:"provider"`
	ProcessingMS int64   `json:"processing_ms"`
}

// StoredResult represents a result stored in the database
type StoredResult struct {
	ID              string
	RequestJSON     []byte
	ResponseJSON    []byte
	Topic           string
	TopicSource     string
	TopicConfidence float64
	CreatedAt       time.Time
}

