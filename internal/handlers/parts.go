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
	cats, _ := app.Categories.List()
	mods, _ := app.Models.List()
	data := map[string]any{
		"Categories": cats,
		"Models":     mods,
	}
	app.render(w, r, "part_new.html", data)
}

func (app *App) PartCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	qty, _ := strconv.Atoi(r.FormValue("quantity"))
	purchasePrice, _ := strconv.ParseFloat(r.FormValue("purchase_price"), 64)
	sellPrice, _ := strconv.ParseFloat(r.FormValue("sell_price"), 64)

	var catID *int64
	if c, err := strconv.ParseInt(r.FormValue("category_id"), 10, 64); err == nil && c > 0 {
		catID = &c
	}

	part := &models.Part{
		Name:          r.FormValue("name"),
		SKU:           r.FormValue("sku"),
		CategoryID:    catID,
		Quantity:      qty,
		PurchasePrice: purchasePrice,
		SellPrice:     sellPrice,
	}

	if part.Name == "" {
		cats, _ := app.Categories.List()
		mods, _ := app.Models.List()
		data := map[string]any{
			"Categories": cats,
			"Models":     mods,
			"Error":      "Название не может быть пустым",
		}
		app.render(w, r, "part_new.html", data)
		return
	}

	r.ParseForm()
	modelStrs := r.Form["models"]
	var modelIDs []int64
	for _, ms := range modelStrs {
		if id, err := strconv.ParseInt(ms, 10, 64); err == nil {
			modelIDs = append(modelIDs, id)
		}
	}

	if _, err := app.Parts.Create(part, modelIDs); err != nil {
		http.Error(w, "Ошибка создания запчасти: "+err.Error(), http.StatusInternalServerError)
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
	
	cats, _ := app.Categories.List()
	mods, _ := app.Models.List()
	
	// Create map of selected models for the template
	selectedModels := make(map[int64]bool)
	for _, m := range part.Models {
		selectedModels[m.ID] = true
	}

	data := map[string]any{
		"Part":           part,
		"Categories":     cats,
		"Models":         mods,
		"SelectedModels": selectedModels,
	}

	app.render(w, r, "part_edit.html", data)
}

func (app *App) PartUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	qty, _ := strconv.Atoi(r.FormValue("quantity"))
	purchasePrice, _ := strconv.ParseFloat(r.FormValue("purchase_price"), 64)
	sellPrice, _ := strconv.ParseFloat(r.FormValue("sell_price"), 64)

	var catID *int64
	if c, err := strconv.ParseInt(r.FormValue("category_id"), 10, 64); err == nil && c > 0 {
		catID = &c
	}

	part := &models.Part{
		ID:            id,
		Name:          r.FormValue("name"),
		SKU:           r.FormValue("sku"),
		CategoryID:    catID,
		Quantity:      qty,
		PurchasePrice: purchasePrice,
		SellPrice:     sellPrice,
	}

	r.ParseForm()
	modelStrs := r.Form["models"]
	var modelIDs []int64
	for _, ms := range modelStrs {
		if id, err := strconv.ParseInt(ms, 10, 64); err == nil {
			modelIDs = append(modelIDs, id)
		}
	}

	if err := app.Parts.Update(part, modelIDs); err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}

	app.flash(r, "Запчасть обновлена")
	http.Redirect(w, r, "/parts", http.StatusSeeOther)
}

func (app *App) PartDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	if err := app.Parts.Delete(id); err != nil {
		app.flash(r, "Ошибка удаления: "+err.Error())
	} else {
		app.flash(r, "Запчасть удалена")
	}
	
	http.Redirect(w, r, "/parts", http.StatusSeeOther)
}
