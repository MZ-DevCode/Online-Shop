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
