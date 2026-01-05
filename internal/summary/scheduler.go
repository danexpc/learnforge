package summary

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	service *Service
	cron    *cron.Cron
}

func NewScheduler(service *Service) *Scheduler {
	return &Scheduler{
		service: service,
		cron:    cron.New(cron.WithLocation(time.UTC)),
	}
}

func (s *Scheduler) Start() {
	_, err := s.cron.AddFunc("0 0 * * *", s.runDailySummary)
	if err != nil {
		log.Printf(`{"level":"error","msg":"Failed to schedule daily summary","error":"%v"}`, err)
		return
	}

	s.cron.Start()
	log.Println(`{"level":"info","msg":"Daily summary scheduler started"}`)
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println(`{"level":"info","msg":"Daily summary scheduler stopped"}`)
}

func (s *Scheduler) runDailySummary() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	summary, err := s.service.GenerateDailySummary(ctx, yesterday)
	if err != nil {
		log.Printf(`{"level":"error","msg":"Failed to generate daily summary","error":"%v"}`, err)
		return
	}

	if err := s.service.SendSummaryToSlack(ctx, summary); err != nil {
		log.Printf(`{"level":"error","msg":"Failed to send daily summary to Slack","error":"%v"}`, err)
		return
	}

	log.Printf(`{"level":"info","msg":"Daily summary sent","date":"%s","requests":%d}`, 
		summary.Date.Format("2006-01-02"), summary.TotalRequests)
}

