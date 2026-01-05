package store

import (
	"context"
	"learnforge/internal/domain"
	"testing"
	"time"
)

func TestInMemStore_SaveAndGet(t *testing.T) {
	store := NewInMemStore()
	ctx := context.Background()

	result := &domain.StoredResult{
		ID:              "test-id",
		RequestJSON:     []byte(`{"text":"test"}`),
		ResponseJSON:    []byte(`{"id":"test-id"}`),
		Topic:           "test-topic",
		TopicSource:     "user",
		TopicConfidence: 1.0,
		CreatedAt:       time.Now(),
	}

	// Save
	if err := store.Save(ctx, result); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Get
	retrieved, err := store.Get(ctx, "test-id")
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}

	if retrieved.ID != result.ID {
		t.Errorf("Expected ID %s, got %s", result.ID, retrieved.ID)
	}

	if retrieved.Topic != result.Topic {
		t.Errorf("Expected topic %s, got %s", result.Topic, retrieved.Topic)
	}
}

func TestInMemStore_GetNotFound(t *testing.T) {
	store := NewInMemStore()
	ctx := context.Background()

	_, err := store.Get(ctx, "non-existent")
	if err == nil {
		t.Fatal("Expected error for non-existent ID")
	}

	domainErr, ok := err.(*domain.DomainError)
	if !ok {
		t.Fatalf("Expected DomainError, got %T", err)
	}

	if domainErr.Code != domain.ErrorCodeNotFound {
		t.Errorf("Expected ErrorCodeNotFound, got %s", domainErr.Code)
	}
}

