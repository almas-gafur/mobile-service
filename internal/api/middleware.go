package api

import (
	"net/http"
	"strings"

	"github.com/example/repair-crm/pkg/auth"
)

func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if allowedOrigin == "*" && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			}
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func authMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				writeError(w, http.StatusUnauthorized, "требуется токен авторизации")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeError(w, http.StatusUnauthorized, "требуется Bearer-токен")
				return
			}

			claims, err := jwtManager.Parse(parts[1])
			if err != nil {
				writeError(w, http.StatusUnauthorized, "некорректный токен")
				return
			}

			next.ServeHTTP(w, withAuthContext(r, authContext{
				MasterID:   claims.MasterID,
				WorkshopID: claims.WorkshopID,
				Username:   claims.Username,
			}))
		})
	}
}
