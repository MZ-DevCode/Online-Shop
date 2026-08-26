package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"net/http"
)

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rows, err := database.DB.Query("SELECT id, name, price, stock FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()

		var products []models.Product

		for rows.Next(){
			var p models.Product
		    err := rows.Scan(&p.id, p.&name, p.&price, p.&stock)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		Product := models.Product{
			ID:    id,
			Name:  name,
			Price: price,
			Stock: stock,
		}
	}

	if r.Method == "POST" {

	}
}
