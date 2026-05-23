package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/example/repair-crm/internal/models"
	"github.com/example/repair-crm/internal/repository"
	"github.com/example/repair-crm/internal/service"
	"github.com/go-chi/chi/v5"
)

type TicketHandler struct {
	tickets *service.TicketService
}

func NewTicketHandler(tickets *service.TicketService) *TicketHandler {
	return &TicketHandler{tickets: tickets}
}

func (h *TicketHandler) SubmitApplication(w http.ResponseWriter, r *http.Request) {
	var input models.PublicApplicationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	ticket, err := h.tickets.CreatePublicApplication(r.Context(), input)
	if service.IsValidationError(err) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeInternalError(w, "не удалось отправить заявку", err)
		return
	}

	writeJSON(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := currentAuth(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var input models.CreateTicketInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	ticket, err := h.tickets.Create(r.Context(), authCtx.WorkshopID, input)
	if service.IsValidationError(err) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeInternalError(w, "не удалось создать заявку", err)
		return
	}

	writeJSON(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := currentAuth(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	tickets, err := h.tickets.List(r.Context(), authCtx.WorkshopID, r.URL.Query().Get("status"))
	if service.IsValidationError(err) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeInternalError(w, "не удалось загрузить заявки", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"tickets": tickets})
}

func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := currentAuth(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	ticket, err := h.tickets.Get(r.Context(), authCtx.WorkshopID, id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "заявка не найдена")
		return
	}
	if err != nil {
		writeInternalError(w, "не удалось загрузить заявку", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) Track(w http.ResponseWriter, r *http.Request) {
	ticket, err := h.tickets.GetPublic(r.Context(), chi.URLParam(r, "hash"))
	if service.IsValidationError(err) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "заявка не найдена")
		return
	}
	if err != nil {
		writeInternalError(w, "не удалось загрузить статус ремонта", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	var input models.ReviewInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	ticket, err := h.tickets.AddReview(r.Context(), chi.URLParam(r, "hash"), input)
	if service.IsValidationError(err) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "заявка не найдена")
		return
	}
	if err != nil {
		writeInternalError(w, "не удалось сохранить отзыв", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := currentAuth(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	var input models.UpdateTicketInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	ticket, err := h.tickets.Update(r.Context(), authCtx.WorkshopID, id, input)
	if service.IsValidationError(err) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "заявка не найдена")
		return
	}
	if err != nil {
		writeInternalError(w, "не удалось обновить заявку", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := currentAuth(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	err := h.tickets.Delete(r.Context(), authCtx.WorkshopID, id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "заявка не найдена")
		return
	}
	if err != nil {
		writeInternalError(w, "не удалось удалить заявку", err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func parseIDParam(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "некорректный идентификатор")
		return 0, false
	}

	return id, true
}
