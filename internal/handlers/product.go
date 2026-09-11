package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"html/template"
	"net/http"
	"strconv"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		userUUID, err := getUserUUIDFromSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Неверный формат цены", http.StatusBadRequest)
			return
		}

		stock, err := strconv.Atoi(r.FormValue("stock"))
		if err != nil {
			http.Error(w, "Неверный формат количества", http.StatusBadRequest)
			return
		}

		p := models.Product{
			Name:        r.FormValue("name"),
			Price:       price,
			Description: r.FormValue("description"),
			Stock:       stock,
		}

		_, err = database.DB.Exec("INSERT INTO")

		tmpl, err := template.ParseFiles("/templates/add_product")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}

	case "GET":
	}
}
