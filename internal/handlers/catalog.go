package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func (c *Context) getUserUUIDFromSession() (string, error) {
	cookie, err := c.R.Cookie("session_id")
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

func (c *Context) CatalogHandler() {
	switch c.R.Method {
	case "GET":
		rows, err := database.DB.Query("SELECT id, name, price, stock FROM products")
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product

		for rows.Next() {
			var p models.Product
			err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
			if err != nil {
				c.Error(err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		tmpl, err := template.ParseFiles("templates/catalog.html")
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(c.W, products)
	}
}

func (c *Context) AddToCart() {
	switch c.R.Method {
	case "POST":
		productID := c.R.FormValue("product_id")
		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			c.Redirect("/login")
			return
		}

		var check int
		err = database.DB.QueryRow("SELECT id FROM cart WHERE user_uuid = ? AND product_id = ?", userUUID, productID).Scan(&check)
		if err == nil {
			c.Redirect("/catalog")
			return
		}

		query := "INSERT INTO cart(user_uuid, product_id) VALUES (?, ?)"

		_, err = database.DB.Exec(query, userUUID, productID)
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		c.Redirect("/catalog")
	}
}

func (c *Context) ShowCart() {
	userUUID, err := c.getUserUUIDFromSession()
	if err != nil {
		c.Redirect("/login")
		return
	}

	var rows *sql.Rows
	rows, err = database.DB.Query(`
		SELECT p.id, p.name, p.price, p.stock
		FROM cart cart_alias
		JOIN products p ON cart_alias.product_id = p.id
		WHERE cart_alias.user_uuid = ?
		`, userUUID)

	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.CartItem
	var totalPrice float64

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
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
		c.Error("Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(c.W, pageData)
}

func (c *Context) RemoveFromCart() {

}
