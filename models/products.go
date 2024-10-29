package models

type Product struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	IsDeleted *bool  `json:"is_deleted"`
}