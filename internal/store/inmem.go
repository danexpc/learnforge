package store

import (
	"context"
	"learnforge/internal/domain"
	"sync"
	"time"
)

type InMemStore struct {
	mu      sync.RWMutex
	results map[string]*domain.StoredResult
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		results: make(map[string]*domain.StoredResult),
	}
}

func (s *InMemStore) Save(ctx context.Context, result *domain.StoredResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[result.ID] = result
	return nil
}

func (s *InMemStore) Get(ctx context.Context, id string) (*domain.StoredResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.results[id]
	if !ok {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "result not found", nil)
	}
	return result, nil
}

func (s *InMemStore) GetByTopic(ctx context.Context, topic string, limit int) ([]*domain.StoredResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var results []*domain.StoredResult
	for _, result := range s.results {
		if result.Topic == topic {
			results = append(results, result)
			if len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

func (s *InMemStore) GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.StoredResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var results []*domain.StoredResult
	for _, result := range s.results {
		if !result.CreatedAt.Before(start) && !result.CreatedAt.After(end) {
			results = append(results, result)
		}
	}
	return results, nil
}

func (s *InMemStore) Close() error {
	return nil
}

