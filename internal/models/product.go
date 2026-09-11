package models

type Product struct {
	ID          int     `json:"id"`
	UserUUID    string  `json:"user_uuid"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
}

type CartItem struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Quantity       int     `json:"quantity"`
	TotalItemPrice float64 `json:"total_item_price"`
}

type CartPageData struct {
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
}
