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

		rows, err := tx.QueryContext(ctx, `SELECT c.quantity, p.price FROM cart c JOIN products p ON c.product_id = p.id WHERE c.user_uuid = ?`, userUUID)
		if err != nil {
			c.Error("Ошибка расчета корзины", http.StatusInternalServerError)
			return
		}

		var totalPrice float64
		for rows.Next() {
			var quantity int
			var price float64

			err := rows.Scan(&quantity, &price)
			if err != nil {
				c.Error(err.Error(), http.StatusInternalServerError)
				return
			}

			totalPrice += float64(quantity) * price

			if quantity <= 0 {
				customError := errors.New("Товар с недопустимым количеством")
				c.Error(customError.Error(), http.StatusBadRequest)
				return
			}
		}
		defer rows.Close()

		c.Redirect("/cart")

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
	}

}
