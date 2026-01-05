package service

import (
	"context"
	"learnforge/internal/ai"
	"learnforge/internal/domain"
	"testing"
	"time"
)

// mockAI is a mock AI client for testing
type mockAI struct {
	processFunc func(ctx context.Context, req *ai.ProcessRequest) (*domain.ProcessResponse, error)
}

func (m *mockAI) ProcessText(ctx context.Context, req *ai.ProcessRequest) (*domain.ProcessResponse, error) {
	if m.processFunc != nil {
		return m.processFunc(ctx, req)
	}
	return &domain.ProcessResponse{
		ID:              "test-id",
		Topic:           "test-topic",
		TopicSource:     "inferred",
		TopicConfidence: 0.9,
		Summary:         "Test summary",
		KeyPoints:       []string{"point1", "point2"},
		Flashcards:      []domain.Flashcard{{Q: "Q1", A: "A1"}},
		Quiz:            []domain.QuizItem{{Q: "Q1", Choices: []string{"A", "B"}, Answer: "A"}},
		Meta:            domain.Meta{Model: "test-model", Provider: "test"},
		CreatedAt:       time.Now(),
	}, nil
}

// mockStore is a mock store for testing
type mockStore struct {
	saveFunc func(ctx context.Context, result *domain.StoredResult) error
	getFunc  func(ctx context.Context, id string) (*domain.StoredResult, error)
}

func (m *mockStore) Save(ctx context.Context, result *domain.StoredResult) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, result)
	}
	return nil
}

func (m *mockStore) Get(ctx context.Context, id string) (*domain.StoredResult, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "not found", nil)
}

func (m *mockStore) GetByTopic(ctx context.Context, topic string, limit int) ([]*domain.StoredResult, error) {
	return nil, nil
}

func (m *mockStore) GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.StoredResult, error) {
	return nil, nil
}

func (m *mockStore) Close() error {
	return nil
}

func TestService_ProcessText_Validation(t *testing.T) {
	svc := NewService(&mockStore{}, &mockAI{})

	tests := []struct {
		name    string
		req     *domain.ProcessRequest
		wantErr bool
	}{
		{
			name:    "empty text",
			req:     &domain.ProcessRequest{Text: ""},
			wantErr: true,
		},
		{
			name:    "invalid mode",
			req:     &domain.ProcessRequest{Text: "test", Mode: "invalid"},
			wantErr: true,
		},
		{
			name:    "invalid level",
			req:     &domain.ProcessRequest{Text: "test", Level: stringPtr("invalid")},
			wantErr: true,
		},
		{
			name:    "valid request",
			req:     &domain.ProcessRequest{Text: "test"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := svc.ProcessText(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ProcessText_Defaults(t *testing.T) {
	svc := NewService(&mockStore{}, &mockAI{})

	req := &domain.ProcessRequest{
		Text: "test text",
	}

	ctx := context.Background()
	resp, err := svc.ProcessText(ctx, req)
	if err != nil {
		t.Fatalf("ProcessText() error = %v", err)
	}

	if resp.ID == "" {
		t.Error("Expected ID to be set")
	}

	if resp.Meta.ProcessingMS < 0 {
		t.Error("Expected processing time to be non-negative")
	}
}

func stringPtr(s string) *string {
	return &s
}

