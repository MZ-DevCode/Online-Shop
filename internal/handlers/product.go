package handlers

import (
	"html/template"
	"net/http"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		userUUID, err := getUserUUIDFromSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("/templates/add_product")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}

	case "GET":
	}
}
