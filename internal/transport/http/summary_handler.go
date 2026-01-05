package http

import (
	"encoding/json"
	"net/http"
	"time"

	"learnforge/internal/domain"
	"learnforge/internal/summary"

	"github.com/go-chi/chi/v5"
)

type SummaryHandler struct {
	service *summary.Service
	apiKey  string
}

func NewSummaryHandler(service *summary.Service, apiKey string) *SummaryHandler {
	return &SummaryHandler{
		service: service,
		apiKey:  apiKey,
	}
}

func (h *SummaryHandler) RegisterRoutes(r chi.Router) {
	r.Post("/v1/summary/generate", h.generateSummary)
}

func (h *SummaryHandler) generateSummary(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" || apiKey != h.apiKey {
		h.writeError(w, http.StatusUnauthorized, domain.ErrorCodeInvalidArgument, "invalid or missing API key", nil)
		return
	}

	dateStr := r.URL.Query().Get("date")
	var targetDate time.Time
	var err error

	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			h.writeError(w, http.StatusBadRequest, domain.ErrorCodeInvalidArgument, "invalid date format (use YYYY-MM-DD)", err)
			return
		}
	} else {
		targetDate = time.Now().UTC().AddDate(0, 0, -1)
	}

	ctx := r.Context()
	sum, err := h.service.GenerateDailySummary(ctx, targetDate)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, domain.ErrorCodeInternal, "failed to generate summary", err)
		return
	}

	if err := h.service.SendSummaryToSlack(ctx, sum); err != nil {
		h.writeError(w, http.StatusInternalServerError, domain.ErrorCodeInternal, "failed to send summary to Slack", err)
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Summary generated and sent to Slack",
		"summary": sum,
	})
}

func (h *SummaryHandler) writeError(w http.ResponseWriter, statusCode int, code domain.ErrorCode, message string, err error) {
	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    string(code),
			"message": message,
		},
	}
	h.writeJSON(w, statusCode, errorResponse)
}

func (h *SummaryHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
