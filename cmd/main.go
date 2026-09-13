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

	mux.HandleFunc("/register", handlers.MakeHandler((*handlers.Context).RegisterHandler))
	mux.HandleFunc("/login", handlers.MakeHandler((*handlers.Context).LoginHandler))
	mux.HandleFunc("/", handlers.MakeHandler((*handlers.Context).CatalogHandler))
	mux.HandleFunc("GET /catalog", handlers.MakeHandler((*handlers.Context).CatalogHandler))
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("POST /cart/add", handlers.MakeHandler((*handlers.Context).AddToCart))
	mux.HandleFunc("GET /cart", handlers.MakeHandler((*handlers.Context).ShowCart))
	mux.HandleFunc("/profile", handlers.MakeHandler((*handlers.Context).ProfileHandler))
	mux.HandleFunc("/profile/password", handlers.MakeHandler((*handlers.Context).ChangePasswordHandler))
	mux.HandleFunc("/profile/product/add", handlers.MakeHandler((*handlers.Context).CreateProductHandler))
	mux.HandleFunc("/logout", handlers.MakeHandler((*handlers.Context).LogoutHandler))
	mux.HandleFunc("cart/remove", handlers.MakeHandler((*handlers.Context).RemoveFromCart))

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
