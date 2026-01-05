package store

import (
	"context"
	"learnforge/internal/domain"
	"time"
)

type Store interface {
	Save(ctx context.Context, result *domain.StoredResult) error
	Get(ctx context.Context, id string) (*domain.StoredResult, error)
	GetByTopic(ctx context.Context, topic string, limit int) ([]*domain.StoredResult, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.StoredResult, error)
	Close() error
}

