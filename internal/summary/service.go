package summary

import (
	"context"
	"encoding/json"
	"fmt"
	"learnforge/internal/cache"
	"learnforge/internal/store"
	"learnforge/internal/slack"
	"sort"
	"time"
)

type Service struct {
	store     store.Store
	cache     cache.Cache
	slack     *slack.Client
	slackErr  *slack.Client
}

func NewService(store store.Store, cache cache.Cache, slackSummary, slackError *slack.Client) *Service {
	return &Service{
		store:    store,
		cache:    cache,
		slack:    slackSummary,
		slackErr: slackError,
	}
}

type DailySummary struct {
	Date        time.Time `json:"date"`
	TotalRequests int     `json:"total_requests"`
	Topics      []TopicStats `json:"topics"`
	Errors      int       `json:"errors"`
}

type TopicStats struct {
	Topic string `json:"topic"`
	Count int    `json:"count"`
}

func (s *Service) GenerateDailySummary(ctx context.Context, date time.Time) (*DailySummary, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	key := fmt.Sprintf("summary:%s", startOfDay.Format("2006-01-02"))
	
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var summary DailySummary
		if err := json.Unmarshal([]byte(cached), &summary); err == nil {
			return &summary, nil
		}
	}

	results, err := s.store.GetByDateRange(ctx, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get results: %w", err)
	}

	summary := &DailySummary{
		Date:          startOfDay,
		TotalRequests: len(results),
		Topics:        make([]TopicStats, 0),
	}

	topicMap := make(map[string]int)
	for _, result := range results {
		topicMap[result.Topic]++
	}

	for topic, count := range topicMap {
		summary.Topics = append(summary.Topics, TopicStats{
			Topic: topic,
			Count: count,
		})
	}

	sort.Slice(summary.Topics, func(i, j int) bool {
		return summary.Topics[i].Count > summary.Topics[j].Count
	})

	summaryJSON, _ := json.Marshal(summary)
	s.cache.Set(ctx, key, string(summaryJSON), 7*24*time.Hour)

	return summary, nil
}

func (s *Service) SendSummaryToSlack(ctx context.Context, summary *DailySummary) error {
	content := fmt.Sprintf(
		"📊 *Daily Summary for %s*\n\n"+
			"• Total Requests: %d\n"+
			"• Topics Processed: %d\n",
		summary.Date.Format("January 2, 2006"),
		summary.TotalRequests,
		len(summary.Topics),
	)

	if len(summary.Topics) > 0 {
		content += "\n*Top Topics:*\n"
		for i, topic := range summary.Topics {
			if i >= 5 {
				break
			}
			content += fmt.Sprintf("• %s: %d\n", topic.Topic, topic.Count)
		}
	}

	return s.slack.SendSummary(ctx, "LearnForge Daily Summary", content)
}

func (s *Service) LogError(ctx context.Context, err error, errorContext map[string]string) {
	if s.slackErr != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			s.slackErr.SendError(ctx, err, errorContext)
		}()
	}
}

