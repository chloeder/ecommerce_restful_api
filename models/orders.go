package models

import (
	"database/sql"
	"errors"
)

type Checkout struct {
	Email    string            `json:"email"`
	Address  string            `json:"address"`
	Products []ProductQuantity `json:"products"`
}

type ProductQuantity struct {
	ID       string `json:"id"`
	Quantity int32  `json:"quantity"`
}

type Order struct {
	ID                string  `json:"id"`
	Email             string  `json:"email"`
	Address           string  `json:"address"`
	GrandTotal        int64   `json:"grand_total"`
	Passcode          *string `json:"passcode,omitempty"`
	PaidAt            *string `json:"paid_at,omitempty"`
	PaidBank          *string `json:"paid_bank,omitempty"`
	PaidAccountNumber *string `json:"paid_account_number,omitempty"`
}

type OrderDetail struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Price     int64  `json:"price"`
	Total     int64  `json:"total"`
}

type OrderWithDetail struct {
	Order
	Details []OrderDetail `json:"details"`
}

func CreateOrder(db *sql.DB, order Order, details []OrderDetail) error {
	// Check if the database is connected
	if db == nil {
		return errors.New("no Database Connected")
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Query to insert order
	queryOrder := `INSERT INTO orders (id, email, address, grand_total, passcode) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(queryOrder, order.ID, order.Email, order.Address, order.GrandTotal, order.Passcode)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	// Query to insert order details
	queryDetail := `INSERT INTO order_details (id, order_id, product_id, quantity, price, total) VALUES ($1, $2, $3, $4, $5, $6)`
	for _, detail := range details {
		_, err = tx.Exec(queryDetail, detail.ID, detail.OrderID, detail.ProductID, detail.Quantity, detail.Price, detail.Total)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
