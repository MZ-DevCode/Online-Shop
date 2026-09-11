package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"html/template"
	"log"
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
			UserUUID:    userUUID,
			Name:        r.FormValue("name"),
			Price:       price,
			Description: r.FormValue("description"),
			Stock:       stock,
		}

		_, err = database.DB.Exec("INSERT INTO products(user_uuid, name, price, description, stock) VALUES (?, ?, ?, ?, ?)", p.UserUUID, p.Name, p.Price, p.Description, p.Stock)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			http.Error(w, "Ошибка записи товара", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/catalog", http.StatusSeeOther)

	case "GET":
		tmpl, err := template.ParseFiles("templates/add_product.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}
