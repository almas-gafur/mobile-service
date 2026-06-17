package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/alexedwards/scs/v2"
	"mobile-service/internal/"
	"mobile-service/internal/"
)

type App struct {
	DB       *sql.DB
	Sessions *scs.SessionManager
	Users    *repository.UserRepo
	Orders   *repository.OrderRepo
	Parts    *repository.PartRepo
	tmpls    map[string]*template.Template
}

func NewApp(db *sql.DB, sm *scs.SessionManager) *App {
	app := &App{
		DB:       db,
		Sessions: sm,
		Users:    repository.NewUserRepo(db),
		Orders:   repository.NewOrderRepo(db),
		Parts:    repository.NewPartRepo(db),
	}
	app.loadTemplates()
	return app
}

func (app *App) loadTemplates() {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "..", "..", "ui", "html")

	funcMap := template.FuncMap{
		"statusLabel": func(s models.OrderStatus) string { return s.Label() },
		"allStatuses": repository.AllStatusLabels,
		"printf":      fmt.Sprintf,
		"slice": func(s string, i, j int) string {
			if i < 0 || j > len(s) || i > j {
				return s
			}
			return s[i:j]
		},
	}

	pages := []string{
		"login.html",
		"orders.html",
		"order_new.html",
		"order_detail.html",
		"parts.html",
		"part_new.html",
		"part_edit.html",
	}

	app.tmpls = make(map[string]*template.Template)
	base := filepath.Join(root, "base.html")

	for _, page := range pages {
		t, err := template.New("").Funcs(funcMap).ParseFiles(base, filepath.Join(root, page))
		if err != nil {
			log.Fatalf("parse template %s: %v", page, err)
		}
		app.tmpls[page] = t
	}

	publicPages := []string{
		"landing.html",
		"track.html",
		"request_success.html",
	}
	publicBase := filepath.Join(root, "public_base.html")

	for _, page := range publicPages {
		t, err := template.New("").Funcs(funcMap).ParseFiles(publicBase, filepath.Join(root, page))
		if err != nil {
			log.Fatalf("parse template %s: %v", page, err)
		}
		app.tmpls[page] = t
	}
}

func (app *App) render(w http.ResponseWriter, r *http.Request, name string, data any) {
	type templateData struct {
		Data        any
		CurrentUser *models.User
		Flash       string
	}

	var currentUser *models.User
	if userID := app.Sessions.GetInt64(r.Context(), "userID"); userID != 0 {
		u, _ := app.Users.GetByID(userID)
		currentUser = u
	}

	flash := app.Sessions.PopString(r.Context(), "flash")

	td := templateData{
		Data:        data,
		CurrentUser: currentUser,
		Flash:       flash,
	}

	t, ok := app.tmpls[name]
	if !ok {
		http.Error(w, fmt.Sprintf("шаблон %s не найден", name), http.StatusInternalServerError)
		return
	}
	// Execute "layout" which will be defined in base.html and public_base.html
	if err := t.ExecuteTemplate(w, "layout", td); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}
}

func (app *App) flash(r *http.Request, msg string) {
	app.Sessions.Put(r.Context(), "flash", msg)
}
