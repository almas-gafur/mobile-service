package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/example/repair-crm/internal/api"
	"github.com/example/repair-crm/internal/repository"
	"github.com/example/repair-crm/internal/service"
	"github.com/example/repair-crm/pkg/auth"
)

func main() {
	cfg := loadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := repository.OpenPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	masterRepo := repository.NewMasterRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, 24*time.Hour)

	authService := service.NewAuthService(masterRepo, jwtManager)
	ticketService := service.NewTicketService(ticketRepo, cfg.DefaultWorkshopID)

	router := api.NewRouter(api.Dependencies{
		AuthService:   authService,
		TicketService: ticketService,
		JWTManager:    jwtManager,
		AllowedOrigin: cfg.AllowedOrigin,
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("backend listening on %s", cfg.HTTPAddr)
		errCh <- server.ListenAndServe()
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-stopCh:
		log.Printf("received signal %s, shutting down", sig)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown server: %v", err)
	}
}

type config struct {
	HTTPAddr          string
	DatabaseURL       string
	JWTSecret         string
	AllowedOrigin     string
	DefaultWorkshopID int64
}

func loadConfig() config {
	return config{
		HTTPAddr:          env("HTTP_ADDR", ":8080"),
		DatabaseURL:       env("DATABASE_URL", "postgres://repair:repair@localhost:5432/repair_crm?sslmode=disable"),
		JWTSecret:         env("JWT_SECRET", "dev-secret-change-me"),
		AllowedOrigin:     env("CORS_ALLOWED_ORIGIN", "*"),
		DefaultWorkshopID: envInt64("DEFAULT_WORKSHOP_ID", 1),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func envInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
