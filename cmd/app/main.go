package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"mobile-service/internal/config"
	"mobile-service/internal/database"
	"mobile-service/internal/handlers"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	// Session manager backed by SQLite
	sessionStore := database.NewSessionStore(db)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sessionStore.CleanupExpired(ctx)

	sm := scs.New()
	sm.Store = sessionStore
	sm.Lifetime = 12 * time.Hour
	sm.Cookie.Secure = false // set true behind HTTPS
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode

	app := handlers.NewApp(db, sm)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(handlers.SecurityHeaders)
	r.Use(sm.LoadAndSave)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))

	// Public routes
	r.Get("/", app.LandingPage)
	r.Post("/request", app.SubmitRequest)
	r.Get("/request/success", app.RequestSuccess)
	r.Get("/track", app.TrackStatus)
	r.Post("/track", app.TrackStatusPost)

	r.Get("/login", app.LoginPage)
	r.Post("/login", app.LoginPost)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(app.RequireAuth)

		r.Post("/logout", app.Logout)

		r.Get("/orders", app.OrdersList)
		r.Get("/orders/new", app.OrderNewPage)
		r.Post("/orders/new", app.OrderCreate)
		r.Get("/orders/{id}", app.OrderDetail)
		r.Post("/orders/{id}/status", app.OrderUpdateStatus)
		r.Post("/orders/{id}/parts", app.OrderWriteOffPart)
		r.Post("/orders/{id}/delete", app.OrderDelete)
		r.Get("/orders/{id}/receipt", app.OrderReceipt)
		r.Get("/orders/{id}/warranty", app.OrderWarranty)

		r.Get("/parts", app.PartsList)
		r.Get("/parts/new", app.PartNewPage)
		r.Post("/parts/new", app.PartCreate)
		r.Get("/parts/{id}/edit", app.PartEditPage)
		r.Post("/parts/{id}/edit", app.PartUpdate)
		r.Post("/parts/{id}/delete", app.PartDelete)

		r.Get("/categories", app.CategoriesList)
		r.Post("/categories/new", app.CategoryCreate)
		r.Post("/categories/{id}/delete", app.CategoryDelete)

		r.Get("/models", app.DeviceModelsList)
		r.Post("/models/new", app.DeviceModelCreate)
		r.Post("/models/{id}/delete", app.DeviceModelDelete)
	})

	log.Printf("Сервер запущен на http://0.0.0.0:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
