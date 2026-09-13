package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"WEBSITE/internal/utils"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

func validatePassword(password string) (bool, string) {
	length := len(password)
	const minLength = 8
	const maxLength = 64

	if length < minLength {
		return false, fmt.Sprintf("Пароль слишком короткий. Минимум %d символов.", minLength)
	}

	if length > maxLength {
		return false, fmt.Sprintf("Пароль слишком длинный. Максимум %d символов.", maxLength)
	}
	return true, "Пароль подходит по длине!"
}

func (c *Context) RegisterHandler() {
	if c.R.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			c.Error("Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(c.W, nil)
		return
	}

	if c.R.Method == "POST" {
		u := models.User{
			Username: c.R.FormValue("username"),
			Name:     c.R.FormValue("name"),
			Password: c.R.FormValue("password"),
			UUID:     utils.GenerateUUID(),
		}

		repeatPassword := c.R.FormValue("repeatPassword")
		if u.Password != repeatPassword {
			c.Error("Пароли не совпадают", http.StatusBadRequest)
			return
		}

		if value, errMsg := validatePassword(u.Password); !value {
			c.Error(errMsg, http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			c.Error("Ошибка шифрования пароля", http.StatusInternalServerError)
			return
		}
		u.Password = hashedPassword

		query := "INSERT INTO users (uuid, name, username, password) VALUES (?, ?, ?, ?)"
		_, err = database.DB.ExecContext(c.Ctx, query, u.UUID, u.Name, u.Username, u.Password)
		if err != nil {
			c.Error("Ошибка регистрации", http.StatusBadRequest)
			return
		}

		sessionToken := utils.GenerateUUID()
		_, err = database.DB.ExecContext(c.Ctx, "INSERT INTO sessions (token, user_uuid) VALUES (?, ?)", sessionToken, u.UUID)
		if err != nil {
			c.Error("Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(c.W, &cookie)

		c.Redirect("/catalog")
	}
}

func (c *Context) LoginHandler() {
	if c.R.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			c.Error("Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(c.W, nil)
		return
	}

	if c.R.Method == "POST" {
		username := c.R.FormValue("username")
		password := c.R.FormValue("password")

		var hashedPassword string
		var userUUID string

		query := "SELECT uuid, password FROM users WHERE username = ?"
		err := database.DB.QueryRowContext(c.Ctx, query, username).Scan(&userUUID, &hashedPassword)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.Error("Неверное имя пользователя или пароль", http.StatusUnauthorized)
				return
			}
			c.Error("Ошибка поиска пользователя", http.StatusInternalServerError)
			return
		}

		if !utils.CheckPasswordHash(password, hashedPassword) {
			c.Error("Неверное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}

		sessionToken := utils.GenerateUUID()

		_, err = database.DB.ExecContext(c.Ctx, "INSERT INTO sessions (token, user_uuid) VALUES (?, ?)", sessionToken, userUUID)
		if err != nil {
			c.Error("Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(c.W, &cookie)

		c.Redirect("/catalog")
	}
}

func (c *Context) LogoutHandler() {
	cookie, err := c.R.Cookie("session_id")
	if err != nil {
		c.Error("Unauthorized: ", http.StatusUnauthorized)
		return
	}
	_, err = database.DB.ExecContext(c.Ctx, "DELETE FROM sessions WHERE token = ?", cookie.Value)
	if err != nil {
		c.Error("Ошибка сервера при выходе", http.StatusInternalServerError)
		return
	}

	cookie = &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(c.W, cookie)
	c.Redirect("/login")
}
