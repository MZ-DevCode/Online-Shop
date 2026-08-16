package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func uuid() string {
	chars := "0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, 12)

	for i := 0; i < 12; i++ {
		randomIndex := r.Intn(len(chars))
		result[i] = chars[randomIndex]
	}

	return string(result)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func checkpassword(password string, repeat string) bool {
	if password == repeat {
		return true
	}
	return false
}

func main() {
	var err error
	db, err = sql.Open("mysql", "user:123@tcp(127.0.0.1:3306)/website_db")

	if err != nil {
		fmt.Print("Error opening database: ", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("Database unreachable: ", err)
		return
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/register.html")
	})
	mux.HandleFunc("POST /register", handlerReg)

	fmt.Println("Server started at http://localhost:8080/register")
	http.ListenAndServe(":8080", mux)
}

func handlerReg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("username")
	password := r.FormValue("password")
	repeat := r.FormValue("repeat_password")
	uuid := uuid()

	if name == "" {
		http.Error(w, "Имя пользователя не может быть пустым", http.StatusBadRequest)
		return
	}

	if !checkpassword(password, repeat) {
		fmt.Println("Ошибка: пароли не совпадают")
		http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
		return
	}

	hash, err := HashPassword(password)
	if err != nil {
		fmt.Println("Ошибка хэширования: ", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users (uuid, username, password) VALUES (?, ?, ?)`
	_, err = db.Exec(query, uuid, name, hash)
	if err != nil {
		fmt.Println("Query error: ", err)
		http.Error(w, "Ошибка базы данных (возможно, имя уже занято)", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Пользователь %s успешно зарегистрирован!", name)
}
