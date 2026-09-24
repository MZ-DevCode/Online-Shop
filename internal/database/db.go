package database

import (
	"WEBSITE/internal/utils"
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

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			telegram_id INTEGER DEFAULT 0
		);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Printf("Ошибка создания таблицы users: %v", err)
	}

	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_uuid TEXT NOT NULL,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		stock INTEGER NOT NULL,
		description TEXT NOT NULL,
		FOREIGN KEY (user_uuid) REFERENCES users(uuid)
	);`

	_, err = DB.Exec(createProductsTable)
	if err != nil {
		log.Printf("Ошибка создания таблицы products: %v", err)
	}

	createCartTable := `
		CREATE TABLE IF NOT EXISTS cart(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_uuid TEXT NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 1,
			FOREIGN KEY (user_uuid) REFERENCES users(uuid),
			FOREIGN KEY (product_id) REFERENCES products(id),
			UNIQUE(user_uuid, product_id)
		);`

	_, err = DB.Exec(createCartTable)
	if err != nil {
		log.Printf("Ошибка создания таблицы cart: %v", err)
	}

	createSessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions(
		token TEXT PRIMARY KEY,
		user_uuid TEXT NOT NULL,
		expires_at DATATIME NOT NULL,
		FOREIGN KEY (user_uuid) REFERENCES users(uuid)
	);`

	_, err = DB.Exec(createSessionsTable)
	if err != nil {
		log.Printf("Ошибка создания таблицы sessions: %v", err)
	}

	createWalletTable := `
	CREATE TABLE IF NOT EXISTS wallets(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_uuid TEXT NOT NULL UNIQUE,
		balance INT NOT NULL DEFAULT 1000,
		FOREIGN KEY (user_uuid) REFERENCES users(uuid)
	);`

	_, err = DB.Exec(createWalletTable)
	if err != nil {
		log.Printf("Ошибка создания таблицы sessions: %v", err)
	}

	hashedPassword, _ := utils.HashPassword("user1")

	_, err = DB.Exec("INSERT OR IGNORE INTO users (uuid, name, username, password) VALUES (?, 'Пользователь 1', 'user1', ?)", utils.GenerateUUID(), hashedPassword)
	if err != nil {
		log.Println("Ошибка вставки юзера 1:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO users (uuid, name, username, password) VALUES (?, 'Пользователь 2', 'user2', ?)", utils.GenerateUUID(), hashedPassword)
	if err != nil {
		log.Println("Ошибка вставки юзера 2:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO users (uuid, name, username, password) VALUES (?, 'Пользователь 3', 'user3', ?)", utils.GenerateUUID(), hashedPassword)
	if err != nil {
		log.Println("Ошибка вставки юзера 3:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO users (uuid, name, username, password) VALUES (?, 'Пользователь 4', 'user4', ?)", utils.GenerateUUID(), hashedPassword)
	if err != nil {
		log.Println("Ошибка вставки юзера 4:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO users (uuid, name, username, password) VALUES (?, 'Пользователь 5', 'user5', ?)", utils.GenerateUUID(), hashedPassword)
	if err != nil {
		log.Println("Ошибка вставки юзера 5:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user1'), 'Товар 1', 100.0, 1, 'Описание 1')")
	if err != nil {
		log.Println("Ошибка вставки товара 1:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user1'), 'Товар 2', 200.0, 1, 'Описание 2')")
	if err != nil {
		log.Println("Ошибка вставки товара 2:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user1'), 'Товар 3', 300.0, 1, 'Описание 3')")
	if err != nil {
		log.Println("Ошибка вставки товара 3:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user1'), 'Товар 4', 400.0, 1, 'Описание 4')")
	if err != nil {
		log.Println("Ошибка вставки товара 4:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user2'), 'Товар 5', 150.0, 1, 'Описание 5')")
	if err != nil {
		log.Println("Ошибка вставки товара 5:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user2'), 'Товар 6', 250.0, 1, 'Описание 6')")
	if err != nil {
		log.Println("Ошибка вставки товара 6:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user2'), 'Товар 7', 350.0, 1, 'Описание 7')")
	if err != nil {
		log.Println("Ошибка вставки товара 7:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user2'), 'Товар 8', 450.0, 1, 'Описание 8')")
	if err != nil {
		log.Println("Ошибка вставки товара 8:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user3'), 'Товар 9', 120.0, 1, 'Описание 9')")
	if err != nil {
		log.Println("Ошибка вставки товара 9:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user3'), 'Товар 10', 220.0, 1, 'Описание 10')")
	if err != nil {
		log.Println("Ошибка вставки товара 10:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user3'), 'Товар 11', 320.0, 1, 'Описание 11')")
	if err != nil {
		log.Println("Ошибка вставки товара 11:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user3'), 'Товар 12', 420.0, 1, 'Описание 12')")
	if err != nil {
		log.Println("Ошибка вставки товара 12:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user4'), 'Товар 13', 180.0, 1, 'Описание 13')")
	if err != nil {
		log.Println("Ошибка вставки товара 13:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user4'), 'Товар 14', 280.0, 1, 'Описание 14')")
	if err != nil {
		log.Println("Ошибка вставки товара 14:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user4'), 'Товар 15', 380.0, 1, 'Описание 15')")
	if err != nil {
		log.Println("Ошибка вставки товара 15:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user4'), 'Товар 16', 480.0, 1, 'Описание 16')")
	if err != nil {
		log.Println("Ошибка вставки товара 16:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user5'), 'Товар 17', 190.0, 1, 'Описание 17')")
	if err != nil {
		log.Println("Ошибка вставки товара 17:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user5'), 'Товар 18', 290.0, 1, 'Описание 18')")
	if err != nil {
		log.Println("Ошибка вставки товара 18:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user5'), 'Товар 19', 390.0, 1, 'Описание 19')")
	if err != nil {
		log.Println("Ошибка вставки товара 19:", err)
	}

	_, err = DB.Exec("INSERT OR IGNORE INTO products (user_uuid, name, price, stock, description) VALUES ((SELECT uuid FROM users WHERE username = 'user5'), 'Товар 20', 490.0, 1, 'Описание 20')")
	if err != nil {
		log.Println("Ошибка вставки товара 20:", err)
	}

	_, err = DB.Exec(`
			INSERT OR IGNORE INTO wallets (user_uuid, balance)
			SELECT uuid, 1000 FROM users;
		`)
	if err != nil {
		log.Println("Ошибка создания кошельков для пользователей:", err)
	}

}
