package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"WEBSITE/internal/utils"
	"html/template"
	"log"
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

		var hash string
		err = database.DB.QueryRow("SELECT password FROM users WHERE uuid = ?", userUUID).Scan(&hash)

		if err != nil {
			http.Error(w, "Ошибка пользователя", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		if !utils.CheckPasswordHash(currentPassword, hash) {
			http.Error(w, "Неверный текущий пароль", http.StatusUnauthorized)
			return
		}

		newHash, err := utils.HashPassword(newPassword)
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		_, err = database.DB.Exec("UPDATE users SET password = ? WHERE uuid = ?", newHash, userUUID)
		if err != nil {
			http.Error(w, "Ошибка сохранения", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

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

	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
