package api

import (
	"context"
	"net/http"
)

type contextKey string

const authContextKey contextKey = "auth"

type authContext struct {
	MasterID   int64
	WorkshopID int64
	Username   string
}

func withAuthContext(r *http.Request, auth authContext) *http.Request {
	ctx := context.WithValue(r.Context(), authContextKey, auth)
	return r.WithContext(ctx)
}

func currentAuth(r *http.Request) (authContext, bool) {
	auth, ok := r.Context().Value(authContextKey).(authContext)
	return auth, ok
}
