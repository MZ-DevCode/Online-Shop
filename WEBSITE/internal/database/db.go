package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite", "online_shop.db")
	if err != nil {
		log.Println("Ошибка открытия базы данных: ", err)
		return
	}

	err = DB.Ping()
	if err != nil {
		log.Println("База данных не отвечает: ", err)
	}

	log.Println("Успешное подкючение к базе данных")

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL
		);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Printf("Ошибка создания таблицы users: %v", err)
	}
}
