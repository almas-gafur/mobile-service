package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"mobile-service/internal/"
	"mobile-service/internal/"
)

type ordersPageData struct {
	Orders       []models.Order
	StatusFilter string
	Search       string
	AllStatuses  []struct {
		Value string
		Label string
	}
	StatusCounts map[string]int
}

func (app *App) OrdersList(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	orders, err := app.Orders.List(status, search)
	if err != nil {
		http.Error(w, "Ошибка загрузки заказов", http.StatusInternalServerError)
		return
	}

	counts, err := app.Orders.StatusCounts()
	if err != nil {
		counts = map[string]int{}
	}

	app.render(w, r, "orders.html", ordersPageData{
		Orders:       orders,
		StatusFilter: status,
		Search:       search,
		AllStatuses:  repository.AllStatusLabels(),
		StatusCounts: counts,
	})
}

func (app *App) OrderNewPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "order_new.html", nil)
}

func (app *App) OrderCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	cost, _ := strconv.ParseFloat(r.FormValue("estimated_cost"), 64)

	order := &models.Order{
		ClientName:    r.FormValue("client_name"),
		Phone:         r.FormValue("phone"),
		Device:        r.FormValue("device"),
		Description:   r.FormValue("description"),
		EstimatedCost: cost,
	}

	if order.ClientName == "" || order.Phone == "" || order.Device == "" {
		app.render(w, r, "order_new.html", map[string]string{"Error": "Заполните обязательные поля"})
		return
	}

	id, err := app.Orders.Create(order)
	if err != nil {
		http.Error(w, "Ошибка создания заказа", http.StatusInternalServerError)
		return
	}

	app.flash(r, "Заказ успешно создан")
	http.Redirect(w, r, "/orders/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

type orderDetailData struct {
	Order       *models.Order
	AllStatuses []struct {
		Value string
		Label string
	}
	Parts []models.Part
}

func (app *App) OrderDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	order, err := app.Orders.GetByID(id)
	if err != nil || order == nil {
		http.NotFound(w, r)
		return
	}

	parts, err := app.Parts.List()
	if err != nil {
		parts = nil
	}

	app.render(w, r, "order_detail.html", orderDetailData{
		Order:       order,
		AllStatuses: repository.AllStatusLabels(),
		Parts:       parts,
	})
}

func (app *App) OrderUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	status := models.OrderStatus(r.FormValue("status"))
	if err := app.Orders.UpdateStatus(id, status); err != nil {
		app.flash(r, "Ошибка: "+err.Error())
	} else {
		app.flash(r, "Статус обновлён")
	}

	http.Redirect(w, r, "/orders/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

func (app *App) OrderDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := app.Orders.Delete(orderID); err != nil {
		app.flash(r, "Ошибка при удалении заказа")
	} else {
		app.flash(r, "Заказ успешно удален")
	}

	http.Redirect(w, r, "/orders", http.StatusSeeOther)
}

func (app *App) OrderWriteOffPart(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	partID, _ := strconv.ParseInt(r.FormValue("part_id"), 10, 64)
	qty, _ := strconv.Atoi(r.FormValue("quantity"))

	if partID == 0 || qty <= 0 {
		app.flash(r, "Укажите запчасть и количество")
		http.Redirect(w, r, "/orders/"+strconv.FormatInt(orderID, 10), http.StatusSeeOther)
		return
	}

	if err := app.Parts.WriteOff(orderID, partID, qty); err != nil {
		app.flash(r, "Ошибка списания: "+err.Error())
	} else {
		app.flash(r, "Запчасть списана со склада")
	}

	http.Redirect(w, r, "/orders/"+strconv.FormatInt(orderID, 10), http.StatusSeeOther)
}
