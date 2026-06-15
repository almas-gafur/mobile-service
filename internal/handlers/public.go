package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"repair-crm/internal/models"
)

// LandingPage renders the public home page with the request form.
func (app *App) LandingPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "landing.html", nil)
}

// SubmitRequest handles public form submissions and creates a new order.
func (app *App) SubmitRequest(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	order := &models.Order{
		ClientName:    r.FormValue("client_name"),
		Phone:         r.FormValue("phone"),
		Device:        r.FormValue("device"),
		Description:   r.FormValue("description"),
		EstimatedCost: 0,
		Status:        models.StatusNew,
	}

	if order.ClientName == "" || order.Phone == "" || order.Device == "" {
		app.render(w, r, "landing.html", map[string]string{"Error": "Заполните обязательные поля"})
		return
	}

	// We need to insert this with StatusNew. 
	// The existing OrderRepo.Create uses models.StatusAccepted.
	// So we need a new method or modify Create to accept status.
	// For now, we will add CreateWithStatus to repository.
	id, err := app.Orders.CreateWithStatus(order, models.StatusNew)
	if err != nil {
		http.Error(w, "Ошибка создания заявки", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/request/success?id="+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

// RequestSuccess shows the success page after a request is submitted.
func (app *App) RequestSuccess(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.render(w, r, "request_success.html", map[string]string{"OrderID": id})
}

// TrackStatus renders the track status page.
func (app *App) TrackStatus(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "track.html", nil)
}

// TrackStatusPost handles the status lookup.
func (app *App) TrackStatusPost(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.FormValue("order_id")
	phone := r.FormValue("phone")

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		app.render(w, r, "track.html", map[string]interface{}{
			"Error": "Неверный номер заказа",
		})
		return
	}

	order, err := app.Orders.GetByID(orderID)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	if order == nil || order.Phone != phone {
		app.render(w, r, "track.html", map[string]interface{}{
			"Error": "Заказ не найден или номер телефона не совпадает",
		})
		return
	}

	app.render(w, r, "track.html", map[string]interface{}{
		"Order":       order,
		"StatusLabel": order.Status.Label(),
	})
}
