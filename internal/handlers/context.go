package handlers

import "net/http"

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func MakeHandler(f func(c *Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(&Context{W: w, R: r})
	}
}

func (c *Context) Redirect(direct string) {
	http.Redirect(c.W, c.R, direct, http.StatusSeeOther)
}

func (c *Context) Error(code int, message string) {
	http.Error(c.W, message, code)
}
