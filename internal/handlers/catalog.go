package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"html/template"
	"log"
	"net/http"
)

func getUserUUIDFromSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	var userUUID string
	err = database.DB.QueryRow("SELECT user_uuid FROM sessions WHERE token = ?", cookie.Value).Scan(&userUUID)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}

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
		userUUID, err := getUserUUIDFromSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var check int
		err = database.DB.QueryRow("SELECT id FROM cart WHERE user_uuid = ? AND product_id = ?", userUUID, productID).Scan(&check)
		if err == nil {
			http.Redirect(w, r, "/catalog", http.StatusSeeOther)
			return
		}

		query := "INSERT INTO cart(user_uuid, product_id) VALUES (?, ?)"

		_, err = database.DB.Exec(query, userUUID, productID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	}
}

func ShowCart(w http.ResponseWriter, r *http.Request) {
	userUUID, err := getUserUUIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := database.DB.Query(`
		SELECT p.id, p.name, p.price, p.stock
		FROM cart c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_uuid = ?
		`, userUUID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.CartItem
	var totalPrice float64

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		item := models.CartItem{
			ID:             p.ID,
			Name:           p.Name,
			Price:          p.Price,
			Quantity:       1,
			TotalItemPrice: p.Price,
		}

		totalPrice += item.TotalItemPrice
		items = append(items, item)
	}

	pageData := models.CartPageData{
		Items:      items,
		TotalPrice: totalPrice,
	}

	tmpl, err := template.ParseFiles("templates/cart.html")
	if err != nil {
		log.Printf("Ошибка загрузки шаблона cart.html: %v", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, pageData)

}
