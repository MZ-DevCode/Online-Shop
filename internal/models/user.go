package models

type User struct {
	ID       int    `db:"id"`
	UUID     string `db:"uuid"`
	Username string `db:"username"`
	Name     string `db:"name"`
	Password string `db:"password"`
}
