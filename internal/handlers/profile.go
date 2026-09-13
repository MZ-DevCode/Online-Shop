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
			c.Redirect("/login")
			return
		}

		currentPassword := c.R.FormValue("current_password")
		newPassword := c.R.FormValue("new_password")

		var hash string
		err = database.DB.QueryRow("SELECT password FROM users WHERE uuid = ?", userUUID).Scan(&hash)

		if err != nil {
			c.Error(http.StatusInternalServerError, "Ошибка пользователя")
			log.Printf("Error: %v", err)
			return
		}

		if !utils.CheckPasswordHash(currentPassword, hash) {
			c.Error(http.StatusUnauthorized, "Неверный текущий пароль")
			return
		}

		newHash, err := utils.HashPassword(newPassword)
		if err != nil {
			c.Error(http.StatusInternalServerError, "Ошибка сервера")
			log.Printf("Error: %v", err)
			return
		}

		_, err = database.DB.Exec("UPDATE users SET password = ? WHERE uuid = ?", newHash, userUUID)
		if err != nil {
			c.Error(http.StatusInternalServerError, "Ошибка сохранения")
			log.Printf("Error: %v", err)
			return
		}

		c.Redirect("/profile")

	default:
		c.Redirect("/profile")
		return
	}
}

func (c *Context) ProfileHandler() {
	switch c.R.Method {
	case "GET":
		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			c.Redirect("/login")
			return
		}

		var u models.User

		err = database.DB.QueryRow("SELECT name, username FROM users WHERE uuid = ?", userUUID).Scan(&u.Name, &u.Username)
		if err != nil {
			c.Error(http.StatusInternalServerError, "Ошибка получения данных пользователя")
			return
		}

		tmpl, err := template.ParseFiles("templates/profile.html")
		if err != nil {
			c.Error(http.StatusInternalServerError, "Ошибка загрузки шаблона")
			return
		}
		tmpl.Execute(c.W, u)

	default:
		c.Redirect("/login")
	}
}
