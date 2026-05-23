package api

import (
	"net/http"

	"github.com/example/repair-crm/internal/service"
	"github.com/example/repair-crm/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	AuthService   *service.AuthService
	TicketService *service.TicketService
	JWTManager    *auth.JWTManager
	AllowedOrigin string
}

func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware(deps.AllowedOrigin))

	authHandler := NewAuthHandler(deps.AuthService)
	ticketHandler := NewTicketHandler(deps.TicketService)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)
		r.Post("/applications", ticketHandler.SubmitApplication)
		r.Get("/track/{hash}", ticketHandler.Track)
		r.Post("/track/{hash}/review", ticketHandler.AddReview)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware(deps.JWTManager))
			r.Get("/tickets", ticketHandler.List)
			r.Post("/tickets", ticketHandler.Create)
			r.Get("/tickets/{id}", ticketHandler.Get)
			r.Put("/tickets/{id}", ticketHandler.Update)
			r.Delete("/tickets/{id}", ticketHandler.Delete)
		})
	})

	return r
}
