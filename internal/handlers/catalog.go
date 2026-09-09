package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"html/template"
	"net/http"
)

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := database.DB.Query("SELECT id, name, price, stock FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product

		for rows.Next() {
			var p models.Product
			err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}
		tmpl, err := template.ParseFiles("templates/catalog.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, products)
	}
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		productID := r.FormValue("product_id")
		userID := 1
		query := "INSERT INTO cart(user_id, product_id) VALUES (?, ?)"

		_, err := database.DB.Exec(query, userID, productID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	}
}

func ShowCart(w http.ResponseWriter, r *http.Request) {
	id := 1

	rows, err := database.DB.Query(`
		SELECT p.id, p.name, p.price, p.stock
		FROM cart c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?
		`, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}
	tmpl, err := template.ParseFiles("templates/cart.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, products)

}
