package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (c *Context) CreateProductHandler() {
	switch c.R.Method {
	case "POST":
		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			c.Redirect("/login")
			return
		}

		price, err := strconv.ParseFloat(c.R.FormValue("price"), 64)
		if err != nil {
			c.Error("Неверный формат цены", http.StatusBadRequest)
			return
		}

		stock, err := strconv.Atoi(c.R.FormValue("stock"))
		if err != nil {
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

		_, err = database.DB.Exec("INSERT INTO products(user_uuid, name, price, description, stock) VALUES (?, ?, ?, ?, ?)", p.UserUUID, p.Name, p.Price, p.Description, p.Stock)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			c.Error("Ошибка записи товара", http.StatusInternalServerError)
			return
		}

		c.Redirect("/catalog")

	case "GET":
		tmpl, err := template.ParseFiles("templates/add_product.html")
		if err != nil {
			c.Error("Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(c.W, nil)
	}
}
