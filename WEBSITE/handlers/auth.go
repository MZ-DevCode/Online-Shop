func handlerReg(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		var name string = r.FormValue("username")
		var password string = r.FormValue("password")
		var repeatPassword string = r.FormValue("password")

		if password != repeatPassword {
			http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
			return
		}
	}
