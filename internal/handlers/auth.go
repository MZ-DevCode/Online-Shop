package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"WEBSITE/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	registerTmpl = template.Must(template.ParseFiles("templates/register.html"))
	loginTmpl    = template.Must(template.ParseFiles("templates/login.html"))
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
	switch c.R.Method {
	case "GET":
		if err := registerTmpl.Execute(c.W, nil); err != nil {
			c.Error("Ошибка загрузки шаблона", http.StatusInternalServerError)
		}

	case "POST":
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

		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		tx, err := database.DB.BeginTx(ctx, nil)

		if err != nil {
			c.Error("Ошибка сервера", http.StatusInternalServerError)
			return
		}

		defer func() {
			if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				log.Printf("Ошибка отката транзакции")
			}
		}()

		query := "INSERT INTO users (uuid, name, username, password) VALUES (?, ?, ?, ?)"
		_, err = tx.ExecContext(ctx, query, u.UUID, u.Name, u.Username, u.Password)
		if err != nil {
			c.Error("Пользователь уже существует или ошибка базы", http.StatusBadRequest)
			return
		}
		_, err = tx.ExecContext(ctx, "INSERT INTO wallets (user_uuid, balance) VALUES (?, 1000)", u.UUID)
		if err != nil {
			c.Error("Ошибка создания кошелька", http.StatusInternalServerError)
			return
		}

		sessionToken := utils.GenerateUUID()
		expiresAt := time.Now().Add(7 * 24 * time.Hour)

		_, err = tx.ExecContext(ctx, "INSERT INTO sessions (token, user_uuid, expires_at) VALUES (?, ?, ?)", sessionToken, u.UUID, expiresAt)
		if err != nil {
			c.Error("Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
			c.Error("Ошибка коммита", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   7 * 24 * 60 * 60,
		}
		http.SetCookie(c.W, &cookie)

		c.Redirect("/catalog")

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *Context) LoginHandler() {
	switch c.R.Method {
	case "GET":
		if err := loginTmpl.Execute(c.W, nil); err != nil {
			c.Error("Ошибка загрузки шаблона", http.StatusInternalServerError)
		}

	case "POST":
		username := c.R.FormValue("username")
		password := c.R.FormValue("password")

		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		var hashedPassword string
		var userUUID string

		query := "SELECT uuid, password FROM users WHERE username = ?"
		err := database.DB.QueryRowContext(ctx, query, username).Scan(&userUUID, &hashedPassword)

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

		_, _ = database.DB.ExecContext(ctx, "DELETE FROM sessions WHERE user_uuid = ? AND expires_at <= CURRENT_TIMESTAMP", userUUID)

		sessionToken := utils.GenerateUUID()
		expiresAt := time.Now().Add(7 * 24 * time.Hour)

		_, err = database.DB.ExecContext(ctx, "INSERT INTO sessions (token, user_uuid, expires_at) VALUES (?, ?, ?)", sessionToken, userUUID, expiresAt)
		if err != nil {
			c.Error("Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   7 * 24 * 60 * 60,
		}
		http.SetCookie(c.W, &cookie)

		c.Redirect("/catalog")

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *Context) LogoutHandler() {
	cookie, err := c.R.Cookie("session_id")
	if err != nil {
		c.Error("Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
	defer cancel()

	_, err = database.DB.ExecContext(ctx, "DELETE FROM sessions WHERE token = ?", cookie.Value)
	if err != nil {
		c.Error("Ошибка сервера при выходе", http.StatusInternalServerError)
		return
	}

	expiredCookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(c.W, expiredCookie)
	c.Redirect("/login")
}
