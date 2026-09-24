package handlers

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/models"
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

		if err := c.R.ParseMultipartForm(10 << 20); err != nil {
			c.Error("Слишком большой файл", http.StatusBadRequest)
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

		file, header, err := c.R.FormFile("image")
		if err != nil {
			c.Error("Ошибка получения изображения", http.StatusBadRequest)
			return
		}
		defer file.Close()

		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.Error("Ошибка сервера при создании папки", http.StatusInternalServerError)
			return
		}

		safeFilename := time.Now().Format("20060102150405") + "_" + header.Filename
		dstPath := filepath.Join(uploadDir, safeFilename)

		dst, err := os.Create(dstPath)
		if err != nil {
			c.Error("Ошибка сохранения файла на сервер", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			c.Error("Ошибка записи файла", http.StatusInternalServerError)
			return
		}

		imageURL := "/uploads/" + safeFilename

		p := models.Product{
			UserUUID:    userUUID,
			Name:        c.R.FormValue("name"),
			Price:       price,
			Description: c.R.FormValue("description"),
			Stock:       stock,
			ImageURL:    imageURL,
		}

		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		query := "INSERT INTO products(user_uuid, name, price, description, stock, image_url) VALUES (?, ?, ?, ?, ?, ?)"
		_, err = database.DB.ExecContext(ctx, query, p.UserUUID, p.Name, p.Price, p.Description, p.Stock, p.ImageURL)
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
