package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"html/template"
	"net/http"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		userUUID, err := getUserUUIDFromSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var p models.Product

		models.Product{
			Name:        r.FormValue("name"),
			Price:       r.FormValue("price"),
			Description: r.FormValue("description"),
			Stock:       r.FormValue("stock"),
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
