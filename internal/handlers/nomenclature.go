package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"mobile-service/internal/models"
)

// Categories List
func (app *App) CategoriesList(w http.ResponseWriter, r *http.Request) {
	cats, err := app.Categories.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, r, "categories.html", cats)
}

// Category Create
func (app *App) CategoryCreate(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		app.flash(r, "Название категории обязательно")
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}
	
	c := &models.Category{Name: name}
	if _, err := app.Categories.Create(c); err != nil {
		app.flash(r, "Ошибка создания: "+err.Error())
	} else {
		app.flash(r, "Категория добавлена")
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

// Category Delete
func (app *App) CategoryDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := app.Categories.Delete(id); err != nil {
		app.flash(r, "Ошибка удаления: "+err.Error())
	} else {
		app.flash(r, "Категория удалена")
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

// DeviceModels List
func (app *App) DeviceModelsList(w http.ResponseWriter, r *http.Request) {
	mods, err := app.Models.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, r, "models.html", mods)
}

// DeviceModel Create
func (app *App) DeviceModelCreate(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		app.flash(r, "Название модели обязательно")
		http.Redirect(w, r, "/models", http.StatusSeeOther)
		return
	}
	
	m := &models.DeviceModel{Name: name}
	if _, err := app.Models.Create(m); err != nil {
		app.flash(r, "Ошибка создания: "+err.Error())
	} else {
		app.flash(r, "Модель добавлена")
	}
	http.Redirect(w, r, "/models", http.StatusSeeOther)
}

// DeviceModel Delete
func (app *App) DeviceModelDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := app.Models.Delete(id); err != nil {
		app.flash(r, "Ошибка удаления: "+err.Error())
	} else {
		app.flash(r, "Модель удалена")
	}
	http.Redirect(w, r, "/models", http.StatusSeeOther)
}
