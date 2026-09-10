package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"html/template"
	"net/http"
)

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		userUUID, err := getUserUUIDFromSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")

	default:
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		userUUID, err := getUserUUIDFromSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var u models.User

		err = database.DB.QueryRow("SELECT name, username FROM users WHERE uuid = ?", userUUID).Scan(&u.Name, &u.Username)
		if err != nil {
			http.Error(w, "Ошибка получения данных пользователя", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/profile.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, u)

	case "POST":

	}
}
