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
			return
		}

		userUUID, err := c.getUserUUIDFromSession()
		if err != nil {
			c.Redirect("/login")
			return
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

		var balance int
		balance, err = tx.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE user_uuid = ?", userUUID).Scan(&balance)
		if err != nil {
			c.Error("Ошибка проверки баланса", http.StatusInternalServerError)
		}

		if balance <= int(totalPrice) {
			c.Error("На балансе нет достаточно средств", http.StatusBadRequest)
			return
		}

		_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance - ? WHERE user_uuid = ?", totalPrice, userUUID)
		if err != nil {
			c.Error("Ошибка считывания денежных средств", http.StatusInternalServerError)
			return
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM FROM cart WHERE user_uuid = ?", userUUID)
		if err != nil {
			c.Error("Ошибка при удалении товарар из корзины", http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
			c.Error("Ошибка коммита", http.StatusInternalServerError)
			return
		}

	default:
		c.Error("Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
