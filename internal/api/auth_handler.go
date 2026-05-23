package api

import (
	"encoding/json"
	"net/http"

	"github.com/example/repair-crm/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	result, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if service.IsValidationError(err) {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err != nil {
		writeInternalError(w, "не удалось войти в систему", err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
