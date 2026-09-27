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

	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	err = DB.Ping()
	if err != nil {
		log.Println("База данных не отвечает: ", err)
	}

	log.Println("Успешное подкючение к базе данных")

	CreateTables()
}
