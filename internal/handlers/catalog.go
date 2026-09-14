package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"context"
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	catalogTmpl = template.Must(template.ParseFiles("templates/catalog.html"))
	cartTmpl    = template.Must(template.ParseFiles("templates/cart.html"))
)

func (c *Context) getUserUUIDFromSession() (string, error) {
	cookie, err := c.R.Cookie("session_id")
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
	defer cancel()

	var userUUID string
	err = database.DB.QueryRowContext(ctx, "SELECT user_uuid FROM sessions WHERE token = ?", cookie.Value).Scan(&userUUID)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}

func (c *Context) CatalogHandler() {
	switch c.R.Method {
	case "GET":
		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		rows, err := database.DB.QueryContext(ctx, "SELECT id, name, price, stock FROM products")
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

		if err := rows.Err(); err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		if err := catalogTmpl.Execute(c.W, products); err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
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

		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		tx, err := database.DB.BeginTx(ctx, nil)
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		defer tx.Rollback()

		var stock int
		err = tx.QueryRowContext(ctx, "SELECT stock FROM products WHERE id = ?", productID).Scan(&stock)

		if errors.Is(err, sql.ErrNoRows) {
			c.Error("Товар не найден", http.StatusNotFound)
			return
		}

		if err != nil {
			c.Error("Ошибка базы данных", http.StatusInternalServerError)
			return
		}

		if stock <= 0 {
			c.Error("Товар закончился", http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO cart (user_uuid, product_id, quantity)
			VALUES (?, ?, 1)
			ON CONFLICT(user_uuid, product_id)
			DO UPDATE SET quantity = quantity + 1
		`

		_, err = tx.ExecContext(ctx, query, userUUID, productID)
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.ExecContext(ctx, "UPDATE products SET stock = stock - 1 WHERE id = ?", productID)
		if err != nil {
			c.Error("Ошибка обновления товара", http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
			c.Error("Ошибка коммита", http.StatusInternalServerError)
			return
		}

		c.Redirect("/catalog")

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
		return

	}
}

func (c *Context) ShowCart() {
	userUUID, err := c.getUserUUIDFromSession()
	if err != nil {
		c.Redirect("/login")
		return
	}

	ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
	defer cancel()

	rows, err := database.DB.QueryContext(ctx, `
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

	if err := cartTmpl.Execute(c.W, pageData); err != nil {
		log.Printf("Ошибка загрузки шаблона cart.html: %v", err)
		c.Error("Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
}

func (c *Context) RemoveFromCart() {
	switch c.R.Method {
	case "POST":
		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			c.Redirect("/login")
			return
		}

		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		productID := c.R.FormValue("product_id")

		_, err = database.DB.ExecContext(ctx, "DELETE FROM cart WHERE user_uuid = ? AND product_id = ?", userUUID, productID)
		if err != nil {
			c.Error("Error", http.StatusInternalServerError)
			return
		}

		c.Redirect("/cart")

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
