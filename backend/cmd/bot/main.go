package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/luvvano/football-bot/backend/internal/config"
	"github.com/luvvano/football-bot/backend/internal/db"
	"github.com/luvvano/football-bot/backend/internal/handlers/api"
	"github.com/luvvano/football-bot/backend/internal/handlers/bot"
	twaMiddleware "github.com/luvvano/football-bot/backend/internal/middleware"
	"github.com/luvvano/football-bot/backend/internal/services"
	tele "gopkg.in/telebot.v3"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config
	cfg := config.Load()

	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	// Connect to database
	database, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Connected to database")

	// Initialize services
	userService := services.NewUserService(database.Pool)
	communityService := services.NewCommunityService(database.Pool)
	venueService := services.NewVenueService(database.Pool)
	eventService := services.NewEventService(database.Pool)
	participantService := services.NewParticipantService(database.Pool)

	// Initialize Telegram bot
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	teleBot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Initialize bot handlers
	botHandler := bot.NewBotHandler(
		teleBot,
		cfg.TelegramWebApp,
		userService,
		communityService,
		venueService,
		eventService,
		participantService,
	)
	botHandler.RegisterHandlers()

	// Initialize HTTP router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize API handlers
	authHandler := api.NewAuthHandler(cfg.TelegramToken, 24*time.Hour)
	eventsHandler := api.NewEventsHandler(eventService, participantService, userService)
	venuesHandler := api.NewVenuesHandler(venueService, communityService, userService)
	participantsHandler := api.NewParticipantsHandler(participantService)
	communitiesHandler := api.NewCommunitiesHandler(communityService)

	// TWA auth middleware
	twaAuth := twaMiddleware.TWAAuthMiddleware(cfg.TelegramToken, 24*time.Hour)

	// Routes
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/api/auth/validate", authHandler.Validate)

	// Communities
	r.Get("/api/communities/{id}", communitiesHandler.GetCommunity)
	r.Get("/api/communities/{id}/events", eventsHandler.ListEvents)
	r.Get("/api/communities/{id}/venues", venuesHandler.ListVenues)
	r.With(twaAuth).Post("/api/communities/{id}/venues", venuesHandler.CreateVenue)

	// Events
	r.Get("/api/events/{id}", eventsHandler.GetEvent)
	r.Get("/api/events/{id}/participants", participantsHandler.ListParticipants)
	r.With(twaAuth).Post("/api/events/{id}/join", eventsHandler.Join)
	r.With(twaAuth).Delete("/api/events/{id}/leave", eventsHandler.Leave)

	// Start HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.APIPort),
		Handler: r,
	}

	go func() {
		log.Printf("API server starting on port %d", cfg.APIPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start bot
	go func() {
		log.Println("Bot starting...")
		teleBot.Start()
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	teleBot.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Shutdown complete")
}
