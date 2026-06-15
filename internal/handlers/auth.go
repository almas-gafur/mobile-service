package handlers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (app *App) LoginPage(w http.ResponseWriter, r *http.Request) {
	if app.Sessions.GetInt64(r.Context(), "userID") != 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.render(w, r, "login.html", nil)
}

func (app *App) LoginPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := app.Users.GetByUsername(username)
	if err != nil || user == nil {
		app.render(w, r, "login.html", map[string]string{"Error": "Неверный логин или пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		app.render(w, r, "login.html", map[string]string{"Error": "Неверный логин или пароль"})
		return
	}

	if err := app.Sessions.RenewToken(r.Context()); err != nil {
		http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
		return
	}

	app.Sessions.Put(r.Context(), "userID", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App) Logout(w http.ResponseWriter, r *http.Request) {
	if err := app.Sessions.Destroy(r.Context()); err != nil {
		http.Error(w, "Ошибка выхода", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
