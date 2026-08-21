package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/utils"
	"html/template"
	"log"
	"net/http"
)

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

		query := "INSERT INTO users (username, password) VALUES (?, ?)"
		_, err = database.DB.Exec(query, username, hashedPassword)
		if err != nil {
			http.Error(w, "Ошибка регистрации", http.StatusBadRequest)
			return
		}
	}
}
