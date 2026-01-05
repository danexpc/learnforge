package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"learnforge/internal/ai"
	"learnforge/internal/domain"
	"learnforge/internal/store"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	store  store.Store
	aiClient ai.Client
}

func NewService(store store.Store, aiClient ai.Client) *Service {
	return &Service{
		store:    store,
		aiClient: aiClient,
	}
}

func (s *Service) ProcessText(ctx context.Context, req *domain.ProcessRequest) (*domain.ProcessResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	if req.IdempotencyKey != nil && *req.IdempotencyKey != "" {
		if existing, err := s.getByIdempotencyKey(ctx, *req.IdempotencyKey); err == nil {
			return existing, nil
		}
	}

	mode := req.Mode
	if mode == "" {
		mode = "lesson"
	}
	language := req.Language
	if language == "" {
		language = "en"
	}

	aiReq := &ai.ProcessRequest{
		Text:     req.Text,
		Mode:     mode,
		Topic:    req.Topic,
		Level:    req.Level,
		Language: language,
	}

	aiCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	startTime := time.Now()
	response, err := s.aiClient.ProcessText(aiCtx, aiReq)
	if err != nil {
		return nil, err
	}
	processingTime := time.Since(startTime).Milliseconds()

	response.Meta.ProcessingMS = processingTime
	response.ID = s.generateID(req)

	if req.Topic != nil && *req.Topic != "" {
		response.Topic = *req.Topic
		response.TopicSource = "user"
		response.TopicConfidence = 1.0
	} else {
		response.TopicSource = "inferred"
	}

	if err := s.saveResult(ctx, req, response); err != nil {
		// Log error but don't fail the request
	}

	return response, nil
}

func (s *Service) GetResult(ctx context.Context, id string) (*domain.ProcessResponse, error) {
	stored, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	var response domain.ProcessResponse
	if err := json.Unmarshal(stored.ResponseJSON, &response); err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeInternal, "failed to unmarshal stored result", err)
	}

	return &response, nil
}

func (s *Service) validateRequest(req *domain.ProcessRequest) error {
	if req.Text == "" {
		return domain.NewDomainError(domain.ErrorCodeInvalidArgument, "text is required", nil)
	}

	if !domain.ValidateMode(req.Mode) {
		return domain.NewDomainError(domain.ErrorCodeInvalidArgument, "mode must be one of: lesson, flashcards, quiz", nil)
	}

	if !domain.ValidateLevel(req.Level) {
		return domain.NewDomainError(domain.ErrorCodeInvalidArgument, "level must be one of: beginner, intermediate, advanced", nil)
	}

	return nil
}

func (s *Service) generateID(req *domain.ProcessRequest) string {
	if req.IdempotencyKey != nil && *req.IdempotencyKey != "" {
		hash := sha256.Sum256([]byte(*req.IdempotencyKey))
		return hex.EncodeToString(hash[:])[:16]
	}
	return uuid.New().String()
}

func (s *Service) getByIdempotencyKey(ctx context.Context, key string) (*domain.ProcessResponse, error) {
	hash := sha256.Sum256([]byte(key))
	id := hex.EncodeToString(hash[:])[:16]
	return s.GetResult(ctx, id)
}

func (s *Service) saveResult(ctx context.Context, req *domain.ProcessRequest, resp *domain.ProcessResponse) error {
	requestJSON, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	responseJSON, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	stored := &domain.StoredResult{
		ID:              resp.ID,
		RequestJSON:     requestJSON,
		ResponseJSON:    responseJSON,
		Topic:           resp.Topic,
		TopicSource:     resp.TopicSource,
		TopicConfidence: resp.TopicConfidence,
		CreatedAt:       resp.CreatedAt,
	}

	return s.store.Save(ctx, stored)
}

