package main

import (
	"WEBSITE/handlers"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.RegisterHandler)
	http.ListenAndServe(":8080", mux)
}
