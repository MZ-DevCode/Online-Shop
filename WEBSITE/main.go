package main

import (
	"html/template"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlerReg)
	http.ListenAndServe(":8080", mux)
}

func handlerReg(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		template.ParseFiles("/templates/register.html")
		var name string = r.FormValue("username")
		var password string = r.FormValue("password")
		var repeatPassword string = r.FormValue("password")

		if password != repeatPassword {
			http.Error(w, "Пароли не совпадают", http.StatusBadRequst)
			return
		}
	}
}
