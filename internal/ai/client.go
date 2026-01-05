package ai

import (
	"context"
	"learnforge/internal/domain"
)

type Client interface {
	ProcessText(ctx context.Context, req *ProcessRequest) (*domain.ProcessResponse, error)
}

type ProcessRequest struct {
	Text     string
	Mode     string // lesson, flashcards, quiz
	Topic    *string
	Level    *string
	Language string
}

