package handlers

import (
	"WEBSITE/internal/database"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"
)

func (c *Context) CheckoutHandler() {
	switch c.R.Method {
	case "POST":
		if err := cartTmpl.Execute(c.W, data); err != nil {
			c.Error("Ошибка отображения")
		}

		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			c.Redirect("/login")
		}

		ctx, cancel := context.WithTimeout(c.Ctx, 3*time.Second)
		defer cancel()

		tx, err := database.DB.BeginTx(ctx, nil)
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		defer func() {
			if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				log.Printf("Ошибка отката транзакции: %v", err)
			}
		}()

		rows, err := tx.QueryContext(ctx, "SELECT с.quantity, p.price FROM cart c JOIN products p ON c.product_id = p.id WHERE id = user_uuid", userUUID)
		for rows.Next() {

		}

		c.Redirect("/cart")

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
	}

}
