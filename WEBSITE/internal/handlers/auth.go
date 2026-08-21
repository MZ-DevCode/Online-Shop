package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"WEBSITE/internal/utils"
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
		username := r.FormValue("username")
		password := r.FormValue("password")
		repeatPassword := r.FormValue("repeatPassword")
		if password != repeatPassword {
			http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Ошибка шифрования пароля", http.StatusInternalServerError)
			return
		}

		user := models.User{
			Username: username,
			Password: hashedPassword,
		}

		query := "INSERT INTO users (username, password) VALUES (?, ?)"
		_, err = database.DB.Exec(query, user.Username, user.Password)
		if err != nil {
			http.Error(w, "Ошибка регистрации", http.StatusBadRequest)
			return
		}
	}
}
