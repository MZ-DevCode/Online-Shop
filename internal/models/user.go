package models

type User struct {
	ID       int     `db:"id"`
	UUID     string  `db:"uuid"`
	Username string  `db:"username"`
	Name     string  `db:"name"`
	Balance  float64 `db:"balance"`
	Password string  `db:"password"`
}
