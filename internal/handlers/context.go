package handlers

import (
	"context"
	"net/http"
)

type Context struct {
	W   http.ResponseWriter
	R   *http.Request
	Ctx context.Context // 1. Добавляем поле сюда
}

func MakeHandler(f func(c *Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(&Context{
			W:   w,
			R:   r,
			Ctx: r.Context(), // 2. Инициализируем один раз для всех запросов
		})
	}
}

func (c *Context) Redirect(direct string) {
	http.Redirect(c.W, c.R, direct, http.StatusSeeOther)
}

func (c *Context) Error(message string, code int) {
	http.Error(c.W, message, code)
}
