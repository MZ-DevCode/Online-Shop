package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"html/template"
	"net/http"
)

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rows, err := database.DB.Query("SELECT id, name, price, stock FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product

		for rows.Next(){
			var p models.Product
		    err := rows.Scan(&p.id, &p.name, &p.price, &p.stock)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)

		tmpl, err := template.ParseFiles("templates/catalog.html")
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, products)
		}

	if r.Method == "POST" {

	}
}
