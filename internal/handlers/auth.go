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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		u := models.User{
			Username: r.FormValue("username"),
			Name:     r.FormValue("name"),
			Password: r.FormValue("password"),
			UUID:     utils.GenerateUUID(),
		}

		repeatPassword := r.FormValue("repeatPassword")
		if u.Password != repeatPassword {
			http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
			return
		}

		if value, errMsg := validatePassword(u.Password); !value {
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			http.Error(w, "Ошибка шифрования пароля", http.StatusInternalServerError)
			return
		}
		u.Password = hashedPassword

		query := "INSERT INTO users (uuid, name, username, password) VALUES (?, ?, ?, ?)"
		_, err = database.DB.Exec(query, u.UUID, u.Name, u.Username, u.Password)
		if err != nil {
			http.Error(w, "Ошибка регистрации", http.StatusBadRequest)
			return
		}

		sessionToken := utils.GenerateUUID()
		_, err = database.DB.Exec("INSERT INTO sessions (token, user_uuid) VALUES (?, ?)", sessionToken, u.UUID)
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var hashedPassword string
		var userUUID string

		query := "SELECT uuid, password FROM users WHERE username = ?"
		err := database.DB.QueryRow(query, username).Scan(&userUUID, &hashedPassword)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Ошибка поиска пользователя", http.StatusInternalServerError)
			return
		}

		if !utils.CheckPasswordHash(password, hashedPassword) {
			http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}

		sessionToken := utils.GenerateUUID()

		_, err = database.DB.Exec("INSERT INTO sessions (token, user_uuid) VALUES (?, ?)", sessionToken, userUUID)
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	}
}
