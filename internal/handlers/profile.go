package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"WEBSITE/internal/utils"
	"context"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	profileTmpl = template.Must(template.ParseFiles("templates/profile.html"))
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

		if value, errMsg := validatePassword(newPassword); !value {
			c.Error(errMsg, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		var hash string
		err = database.DB.QueryRowContext(ctx, "SELECT password FROM users WHERE uuid = ?", userUUID).Scan(&hash)
		if err != nil {
			c.Error("Ошибка пользователя", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		if !utils.CheckPasswordHash(currentPassword, hash) {
			c.Error("Неверный текущий пароль", http.StatusUnauthorized)
			return
		}

		newHash, err := utils.HashPassword(newPassword)
		if err != nil {
			c.Error("Ошибка сервера", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		_, err = database.DB.ExecContext(ctx, "UPDATE users SET password = ? WHERE uuid = ?", newHash, userUUID)
		if err != nil {
			c.Error("Ошибка сохранения", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		c.LogoutHandler()

	default:
		c.Redirect("/profile")
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

		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		type userResult struct {
			user models.User
			err  error
		}

		type balanceResult struct {
			balance float64
			err     error
		}

		userChan := make(chan userResult, 1)
		balanceChan := make(chan balanceResult, 1)

		go func() {
			var u models.User
			err := database.DB.QueryRowContext(ctx, "SELECT uuid, name, username FROM users WHERE uuid = ?", userUUID).Scan(&u.UUID, &u.Name, &u.Username)
			userChan <- userResult{
				user: u,
				err:  err,
			}

		}()

		go func() {
			var balance float64
			err := database.DB.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE user_uuid = ?", userUUID).Scan(&balance)
			balanceChan <- balanceResult{
				balance: balance,
				err:     err,
			}

		}()

		uRes := <-userChan
		if uRes.err != nil {
			c.Error("Ошибка получения данных пользователя", http.StatusInternalServerError)
			return
		}

		bRes := <-balanceChan
		if bRes.err != nil {
			c.Error("Ошибка получения баланса пользователя", http.StatusInternalServerError)
			return
		}

		u := uRes.user
		u.Balance = bRes.balance

		if err := profileTmpl.Execute(c.W, u); err != nil {
			c.Error("Ошибка загрузки шаблона", http.StatusInternalServerError)
		}

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
	}
}
