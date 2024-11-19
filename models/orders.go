package models

import (
	"database/sql"
	"errors"
)

// ProductQuantity is representing the product ID and quantity in API request
type ProductQuantity struct {
	ID       string `json:"id" binding:"required"`
	Quantity int32  `json:"quantity" binding:"required"`
}

// Checkout is representing the checkout in API request
type Checkout struct {
	Email    string            `json:"email" binding:"required,email"`
	Address  string            `json:"address" binding:"required"`
	Products []ProductQuantity `json:"products"  binding:"min=1"`
}

// OrderConfirmation is representing the order confirmation in API response
type OrderConfirmation struct {
	Amount        int64  `json:"amount" binding:"required"`
	Bank          string `json:"bank" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required"`
	Passcode      string `json:"passcode"`
}

// Order is representing the order in database
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

// OrderDetail is representing the order detail in database
type OrderDetail struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Price     int64  `json:"price"`
	Total     int64  `json:"total"`
}

// OrderWithDetail is representing the order with detail in API response (without passcode)
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

func SelectOrderById(db *sql.DB, id string) (Order, error) {
	// Check if the database is connected
	if db == nil {
		return Order{}, errors.New("No Database Connected")
	}

	// Query to select order by id
	queryOrder := `SELECT id, email, address, grand_total, passcode, paid_at, paid_bank, paid_account_number FROM orders WHERE id = $1`
	row := db.QueryRow(queryOrder, id)

	// Make variable to store the result
	order := Order{}

	// Scan the result to the variable
	err := row.Scan(&order.ID, &order.Email, &order.Address, &order.GrandTotal, &order.Passcode, &order.PaidAt, &order.PaidBank, &order.PaidAccountNumber)
	if err != nil {
		return Order{}, err
	}

	// Return the result
	return order, nil
}
