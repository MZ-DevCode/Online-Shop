package main

import (
	_ "WEBSITE/docs"
	"WEBSITE/internal/database"
	"WEBSITE/internal/handlers"
	"log"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Интернет-магазин API
// @version 1.0
// @description бэкенд интернет-магазин
// @host localhost:8080
// @BasePath /

func main() {

	database.InitDB()
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/", handlers.CatalogHandler)
	mux.HandleFunc("GET /catalog", handlers.CatalogHandler)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("POST /cart/add", handlers.AddToCart)
	mux.HandleFunc("GET /cart", handlers.ShowCart)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Сервер запущен на http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error to started server:", err)
	}
}
