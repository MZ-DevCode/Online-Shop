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
