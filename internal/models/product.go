package models

type Product struct {
	ID          int     `json:"id"`
	UserUUID    string  `json:"user_uuid"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
}
