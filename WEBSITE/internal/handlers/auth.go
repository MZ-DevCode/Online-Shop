package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"WEBSITE/internal/utils"
	"database/sql"
	"errors"
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
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		repeatPassword := r.FormValue("repeatPassword")
		if u.Password != repeatPassword {
			http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			http.Error(w, "Ошибка шифрования пароля", http.StatusInternalServerError)
			return
		}
		u.Password = hashedPassword

		query := "INSERT INTO users (username, password) VALUES (?, ?)"
		_, err = database.DB.Exec(query, u.Username, u.Password)
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
		query := "SELECT password FROM users WHERE username = $1"
		err := database.DB.QueryRow(query, username).Scan(&hashedPassword)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println("Ошибка поиска пользователя:", err)
				http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
				return
			}

			http.Error(w, "Ошибка поиска пользователя", http.StatusInternalServerError)
			return
		}

		if !utils.CheckPasswordHash(password, hashedPassword) {
			log.Println("Неверный пароль для пользователя:", username)
			http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}

		w.Write([]byte("Успешная авторизация!"))
	}
}
