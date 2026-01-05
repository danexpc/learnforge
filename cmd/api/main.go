package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"learnforge/internal/ai"
	"learnforge/internal/cache"
	"learnforge/internal/config"
	"learnforge/internal/service"
	"learnforge/internal/slack"
	"learnforge/internal/store"
	"learnforge/internal/summary"
	httptransport "learnforge/internal/transport/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.SetFlags(0)

	var st store.Store
	if cfg.Storage == "postgres" {
		if cfg.DatabaseURL == "" {
			log.Fatal("DATABASE_URL is required when STORAGE=postgres")
		}
		st, err = store.NewPostgresStore(cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		log.Println(`{"level":"info","msg":"Using PostgreSQL storage"}`)
	} else {
		st = store.NewInMemStore()
		log.Println(`{"level":"info","msg":"Using in-memory storage"}`)
	}
	defer st.Close()

	if cfg.AIApiKey == "" {
		log.Fatal("AI_API_KEY is required")
	}

	var aiClient ai.Client
	if cfg.AIProvider == "gemini" {
		aiClient = ai.NewGeminiClient(cfg.AIApiKey, cfg.AIModel)
		log.Println(`{"level":"info","msg":"Using Gemini AI provider"}`)
	} else {
		aiClient = ai.NewOpenAIClient(cfg.AIBaseURL, cfg.AIApiKey, cfg.AIModel)
		log.Println(`{"level":"info","msg":"Using OpenAI AI provider"}`)
	}

	svc := service.NewService(st, aiClient)

	var cacheClient cache.Cache
	if cfg.RedisURL != "" {
		redisCache, err := cache.NewRedisCache(cfg.RedisURL)
		if err != nil {
			log.Printf(`{"level":"warn","msg":"Failed to connect to Redis, using in-memory cache","error":"%v"}`, err)
			cacheClient = cache.NewInMemCache()
		} else {
			cacheClient = redisCache
			log.Println(`{"level":"info","msg":"Using Redis cache"}`)
		}
	} else {
		cacheClient = cache.NewInMemCache()
		log.Println(`{"level":"info","msg":"Using in-memory cache"}`)
	}
	defer cacheClient.Close()

	var slackSummary, slackError *slack.Client
	if cfg.SlackWebhookURL != "" {
		slackSummary = slack.NewClient(cfg.SlackWebhookURL)
	}
	if cfg.SlackErrorWebhookURL != "" {
		slackError = slack.NewClient(cfg.SlackErrorWebhookURL)
	}

	summarySvc := summary.NewService(st, cacheClient, slackSummary, slackError)
	summaryScheduler := summary.NewScheduler(summarySvc)
	if slackSummary != nil {
		summaryScheduler.Start()
		defer summaryScheduler.Stop()
	}

	handler := httptransport.NewHandler(svc, summarySvc)

	r := chi.NewRouter()
	r.Use(httptransport.RequestIDMiddleware)
	r.Use(httptransport.LoggingMiddleware)
	r.Use(httptransport.MetricsMiddleware)

	handler.RegisterRoutes(r)

	if cfg.SummaryAPIKey != "" {
		summaryHandler := httptransport.NewSummaryHandler(summarySvc, cfg.SummaryAPIKey)
		summaryHandler.RegisterRoutes(r)
	}

	// Register web UI routes (must be last to catch-all non-API routes)
	handler.RegisterWebRoutes(r)

	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second, // Increased for meme generation
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf(`{"level":"info","msg":"Starting server","port":"%s"}`, cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println(`{"level":"info","msg":"Shutting down server..."}`)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println(`{"level":"info","msg":"Server exited"}`)
}
