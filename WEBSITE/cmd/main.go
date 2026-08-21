package main

import (
	_ "WEBSITE/docs"
	"WEBSITE/internal/database"
	"WEBSITE/internal/handlers"
	"log"
	"net/http"

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
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
