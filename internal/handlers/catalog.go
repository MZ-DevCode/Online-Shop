package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
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
	err = database.DB.QueryRowContext(c.Ctx, "SELECT user_uuid FROM sessions WHERE token = ?", cookie.Value).Scan(&userUUID)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}

func (c *Context) CatalogHandler() {
	switch c.R.Method {
	case "GET":
		rows, err := database.DB.QueryContext(c.Ctx, "SELECT id, name, price, stock FROM products")
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

		tx, err := database.DB.BeginTx(c.Ctx, nil)
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		defer tx.Rollback()

		var stock int
		err = tx.QueryRowContext(c.Ctx, "SELECT stock FROM products WHERE id = ?", productID).Scan(&stock)
		if err != nil || stock <= 0 {
			c.Error("Товар закончился", http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO cart (user_uuid, product_id, quantity)
			VALUES (?, ?, 1)
			ON CONFLICT(user_uuid, product_id)
			DO UPDATE SET quantity = quantity + 1
		`

		_, err = tx.ExecContext(c.Ctx, query, userUUID, productID)
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.ExecContext(c.Ctx, "UPDATE products SET stock = stock - 1 WHERE id = ?", productID)
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
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

	rows, err := database.DB.QueryContext(c.Ctx, `
		SELECT p.id, p.name, p.price, p.stock, c_alias.quantity
		FROM cart c_alias
		JOIN products p ON c_alias.product_id = p.id
		WHERE c_alias.user_uuid = ?
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
		var quantity int
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &quantity)
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		itemPrice := p.Price * float64(quantity)

		item := models.CartItem{
			ID:             p.ID,
			Name:           p.Name,
			Price:          p.Price,
			Quantity:       quantity,
			TotalItemPrice: itemPrice,
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
	switch c.R.Method {
	case "POST":
		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			c.Redirect("/login")
			return
		}

		productID := c.R.FormValue("product_id")

		_, err = database.DB.ExecContext(c.Ctx, "DELETE FROM cart WHERE user_uuid = ? AND product_id = ?", userUUID, productID)
		if err != nil {
			c.Error("Error", http.StatusInternalServerError)
			return
		}

		c.Redirect("/cart")

	default:
		c.Redirect("/cart")
	}
}
