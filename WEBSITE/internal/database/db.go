package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(){}
	var err error

	DB, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/online_shop")
	if err != nil{
		log.Println("Ошибка открытия базы данных: ", err)
		return
	}

	err = DB.Ping()
	if err != nil{
		log.Println("База данных не отвечает: ", err)
	}

	log.Println("Успешное подкючение к базе данных")
}
