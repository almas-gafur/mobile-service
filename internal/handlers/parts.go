package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"mobile-service/internal/models"
)

func (app *App) PartsList(w http.ResponseWriter, r *http.Request) {
	parts, err := app.Parts.List()
	if err != nil {
		http.Error(w, "Ошибка загрузки склада", http.StatusInternalServerError)
		return
	}
	app.render(w, r, "parts.html", parts)
}

func (app *App) PartNewPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "part_new.html", nil)
}

func (app *App) PartCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	qty, _ := strconv.Atoi(r.FormValue("quantity"))
	price, _ := strconv.ParseFloat(r.FormValue("purchase_price"), 64)

	part := &models.Part{
		Name:          r.FormValue("name"),
		Quantity:      qty,
		PurchasePrice: price,
	}

	if part.Name == "" {
		app.render(w, r, "part_new.html", map[string]string{"Error": "Название не может быть пустым"})
		return
	}

	if _, err := app.Parts.Create(part); err != nil {
		http.Error(w, "Ошибка создания запчасти", http.StatusInternalServerError)
		return
	}

	app.flash(r, "Запчасть добавлена на склад")
	http.Redirect(w, r, "/parts", http.StatusSeeOther)
}

func (app *App) PartEditPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	part, err := app.Parts.GetByID(id)
	if err != nil || part == nil {
		http.NotFound(w, r)
		return
	}
	app.render(w, r, "part_edit.html", part)
}

func (app *App) PartUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	qty, _ := strconv.Atoi(r.FormValue("quantity"))
	price, _ := strconv.ParseFloat(r.FormValue("purchase_price"), 64)

	part := &models.Part{
		ID:            id,
		Name:          r.FormValue("name"),
		Quantity:      qty,
		PurchasePrice: price,
	}

	if err := app.Parts.Update(part); err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}

	app.flash(r, "Запчасть обновлена")
	http.Redirect(w, r, "/parts", http.StatusSeeOther)
}
