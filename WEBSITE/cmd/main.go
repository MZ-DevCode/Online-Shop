package main

import (
	"WEBSITE/internal/database"
	"WEBSITE/internal/handlers"
	"log"
	"net/http"
)

func main() {

	database.InitDB()
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/register", handlers.RegisterHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
