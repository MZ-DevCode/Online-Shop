package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"WEBSITE/internal/utils"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// RegisterHandler обрабатывает регистрацию нового пользователя
// @Summary Регистрация пользователя
// @Description Отображает форму регистрации (GET) или создает нового пользователя (POST)
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce text/html
// @Param username formData string true "Имя пользователя"
// @Param password formData string true "Пароль"
// @Param repeatPassword formData string true "Повтор пароля"
// @Success 200 {string} string "HTML страница или успешная регистрация"
// @Failure 400 {string} string "Ошибка валидации или несовпадение паролей"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /register [get]
// @Router /register [post]

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
			log.Println("Ошибка загрузки шаблона:", err)
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		u := models.User{
			Name:     r.FormValue("name"),
			Password: r.FormValue("password"),
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

		query := "INSERT INTO users (name, password) VALUES (?, ?)"
		_, err = database.DB.Exec(query, u.Name, u.Password)
		if err != nil {
			http.Error(w, "Ошибка регистрации", http.StatusBadRequest)
			return
		}
		w.Write([]byte("Регистрация успешна!"))
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println("Ошибка загрузки шаблона:", err)
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

		query := "SELECT password FROM users WHERE username = ?"
		err := database.DB.QueryRow(query, username).Scan(&hashedPassword)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println("Ошибка поиска пользователя:", err)
				http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
				return
			}

			log.Println("Ошибка БД:", err)
			http.Error(w, "Ошибка поиска пользователя", http.StatusInternalServerError)
			return
		}

		if !utils.CheckPasswordHash(password, hashedPassword) {
			log.Println("Неверный пароль для пользователя:", username)
			http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}

		w.Write([]byte("Успешная авторизация!"))
		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	}
}
