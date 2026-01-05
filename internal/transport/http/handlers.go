package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"learnforge/internal/domain"
	"learnforge/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service        *service.Service
	summaryService interface {
		LogError(ctx context.Context, err error, context map[string]string)
	}
}

func NewHandler(service *service.Service, summaryService interface {
	LogError(ctx context.Context, err error, context map[string]string)
}) *Handler {
	return &Handler{
		service:        service,
		summaryService: summaryService,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/v1/process", h.processText)
	r.Get("/v1/process/{id}", h.getResult)
	r.Get("/healthz", h.healthz)
	r.Get("/readyz", h.readyz)
	r.Get("/metrics", h.metrics)
	
	h.RegisterOpenAPIRoutes(r)
}

func (h *Handler) processText(w http.ResponseWriter, r *http.Request) {
	var req domain.ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, domain.ErrorCodeInvalidArgument, "invalid request body", err)
		return
	}

	ctx := r.Context()
	response, err := h.service.ProcessText(ctx, &req)
	if err != nil {
		if h.summaryService != nil {
			requestID := ctx.Value("request_id")
			h.summaryService.LogError(ctx, err, map[string]string{
				"request_id": fmt.Sprintf("%v", requestID),
				"endpoint":   "/v1/process",
			})
		}
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) getResult(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, domain.ErrorCodeInvalidArgument, "id is required", nil)
		return
	}

	ctx := r.Context()
	response, err := h.service.GetResult(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) readyz(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("# Metrics endpoint\n"))
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	domainErr, ok := err.(*domain.DomainError)
	if !ok {
		h.writeError(w, http.StatusInternalServerError, domain.ErrorCodeInternal, "internal server error", err)
		return
	}

	var statusCode int
	switch domainErr.Code {
	case domain.ErrorCodeInvalidArgument:
		statusCode = http.StatusBadRequest
	case domain.ErrorCodeNotFound:
		statusCode = http.StatusNotFound
	case domain.ErrorCodeUpstreamTimeout:
		statusCode = http.StatusGatewayTimeout
	case domain.ErrorCodeUpstreamError:
		statusCode = http.StatusBadGateway
	default:
		statusCode = http.StatusInternalServerError
	}

	h.writeError(w, statusCode, domainErr.Code, domainErr.Message, domainErr.Err)
}

func (h *Handler) writeError(w http.ResponseWriter, statusCode int, code domain.ErrorCode, message string, err error) {
	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    string(code),
			"message": message,
		},
	}

	h.writeJSON(w, statusCode, errorResponse)
}

func (h *Handler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
