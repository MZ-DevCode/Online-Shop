package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	addProductTmpl = template.Must(template.ParseFiles("templates/add_product.html"))
)

func (c *Context) CreateProductHandler() {
	switch c.R.Method {
	case "GET":
		if err := addProductTmpl.Execute(c.W, nil); err != nil {
			c.Error("Ошибка загрузки шаблона", http.StatusInternalServerError)
		}

	case "POST":
		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			c.Redirect("/login")
			return
		}

		price, err := strconv.ParseFloat(c.R.FormValue("price"), 64)
		if err != nil || price < 0 {
			c.Error("Неверный формат цены", http.StatusBadRequest)
			return
		}

		stock, err := strconv.Atoi(c.R.FormValue("stock"))
		if err != nil || stock < 0 {
			c.Error("Неверный формат количества", http.StatusBadRequest)
			return
		}

		p := models.Product{
			UserUUID:    userUUID,
			Name:        c.R.FormValue("name"),
			Price:       price,
			Description: c.R.FormValue("description"),
			Stock:       stock,
		}

		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		query := "INSERT INTO products(user_uuid, name, price, description, stock) VALUES (?, ?, ?, ?, ?)"
		_, err = database.DB.ExecContext(ctx, query, p.UserUUID, p.Name, p.Price, p.Description, p.Stock)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			c.Error("Ошибка записи товара", http.StatusInternalServerError)
			return
		}

		c.Redirect("/catalog")

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
	}
}
