package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"WEBSITE/internal/utils"
	"html/template"
	"log"
	"net/http"
)

func (c *Context) ChangePasswordHandler() {
	switch c.R.Method {
	case "POST":
		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			http.Redirect(c.W, c.R, "/login", http.StatusSeeOther)
			return
		}

		currentPassword := c.R.FormValue("current_password")
		newPassword := c.R.FormValue("new_password")

		var hash string
		err = database.DB.QueryRow("SELECT password FROM users WHERE uuid = ?", userUUID).Scan(&hash)

		if err != nil {
			http.Error(c.W, "Ошибка пользователя", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		if !utils.CheckPasswordHash(currentPassword, hash) {
			http.Error(c.W, "Неверный текущий пароль", http.StatusUnauthorized)
			return
		}

		newHash, err := utils.HashPassword(newPassword)
		if err != nil {
			http.Error(c.W, "Ошибка сервера", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		_, err = database.DB.Exec("UPDATE users SET password = ? WHERE uuid = ?", newHash, userUUID)
		if err != nil {
			http.Error(c.W, "Ошибка сохранения", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		http.Redirect(c.W, c.R, "/profile", http.StatusSeeOther)

	default:
		http.Redirect(c.W, c.R, "/profile", http.StatusSeeOther)
		return
	}
}

func (c *Context) ProfileHandler() {
	switch c.R.Method {
	case "GET":
		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			http.Redirect(c.W, c.R, "/login", http.StatusSeeOther)
			return
		}

		var u models.User

		err = database.DB.QueryRow("SELECT name, username FROM users WHERE uuid = ?", userUUID).Scan(&u.Name, &u.Username)
		if err != nil {
			http.Error(c.W, "Ошибка получения данных пользователя", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/profile.html")
		if err != nil {
			http.Error(c.W, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(c.W, u)

	default:
		http.Redirect(c.W, c.R, "/login", http.StatusSeeOther)
	}
}
