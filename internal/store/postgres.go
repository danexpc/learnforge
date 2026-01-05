package store

import (
	"context"
	"database/sql"
	"learnforge/internal/domain"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(databaseURL string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &PostgresStore{db: db}
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *PostgresStore) Save(ctx context.Context, result *domain.StoredResult) error {
	query := `
		INSERT INTO processed_results (id, request_json, response_json, topic, topic_source, topic_confidence, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			request_json = EXCLUDED.request_json,
			response_json = EXCLUDED.response_json,
			topic = EXCLUDED.topic,
			topic_source = EXCLUDED.topic_source,
			topic_confidence = EXCLUDED.topic_confidence
	`

	_, err := s.db.ExecContext(ctx, query,
		result.ID,
		result.RequestJSON,
		result.ResponseJSON,
		result.Topic,
		result.TopicSource,
		result.TopicConfidence,
		result.CreatedAt,
	)
	return err
}

func (s *PostgresStore) Get(ctx context.Context, id string) (*domain.StoredResult, error) {
	query := `
		SELECT id, request_json, response_json, topic, topic_source, topic_confidence, created_at
		FROM processed_results
		WHERE id = $1
	`

	var result domain.StoredResult
	var createdAt time.Time
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.RequestJSON,
		&result.ResponseJSON,
		&result.Topic,
		&result.TopicSource,
		&result.TopicConfidence,
		&createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "result not found", nil)
	}
	if err != nil {
		return nil, err
	}

	result.CreatedAt = createdAt
	return &result, nil
}

func (s *PostgresStore) GetByTopic(ctx context.Context, topic string, limit int) ([]*domain.StoredResult, error) {
	query := `
		SELECT id, request_json, response_json, topic, topic_source, topic_confidence, created_at
		FROM processed_results
		WHERE topic = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.QueryContext(ctx, query, topic, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.StoredResult
	for rows.Next() {
		var result domain.StoredResult
		var createdAt time.Time
		if err := rows.Scan(
			&result.ID,
			&result.RequestJSON,
			&result.ResponseJSON,
			&result.Topic,
			&result.TopicSource,
			&result.TopicConfidence,
			&createdAt,
		); err != nil {
			return nil, err
		}
		result.CreatedAt = createdAt
		results = append(results, &result)
	}

	return results, rows.Err()
}

func (s *PostgresStore) GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.StoredResult, error) {
	query := `
		SELECT id, request_json, response_json, topic, topic_source, topic_confidence, created_at
		FROM processed_results
		WHERE created_at >= $1 AND created_at < $2
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.StoredResult
	for rows.Next() {
		var result domain.StoredResult
		var createdAt time.Time
		if err := rows.Scan(
			&result.ID,
			&result.RequestJSON,
			&result.ResponseJSON,
			&result.Topic,
			&result.TopicSource,
			&result.TopicConfidence,
			&createdAt,
		); err != nil {
			return nil, err
		}
		result.CreatedAt = createdAt
		results = append(results, &result)
	}

	return results, rows.Err()
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

