package handlers

import (
	"html/template"
	"net/http"
)

func (c *Context) CheckoutHandler() {
	switch c.R.Method {
	case "POST":
		if err := cartTmpl.Execute(c.W, data); err != nil{
			c.Error("Ошибка отображения")
		}

			c.Redirect("/cart")
	}

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
	}
}
