package migration

import (
	"database/sql"
	"log"
)

func Migration(db *sql.DB) error {
	if db == nil {
		log.Fatal("Database not connected")
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
		    id VARCHAR(36) PRIMARY KEY,
		    name VARCHAR(255) NOT NULL,
		    price INT NOT NULL,
			is_deleted BOOLEAN DEFAULT FALSE);

		CREATE TABLE IF NOT EXISTS orders (
		    id VARCHAR(36) PRIMARY KEY,
		    email VARCHAR(255) NOT NULL,
		    address VARCHAR(255) NOT NULL,
		    passcode VARCHAR,
		    paid_at TIMESTAMP,
		    paid_bank VARCHAR(255),
		    paid_account_number VARCHAR(255),
		    grand_total BIGINT NOT NULL);

		CREATE TABLE IF NOT EXISTS order_details (
		    id VARCHAR(36) PRIMARY KEY,
		    order_id VARCHAR(36) NOT NULL,
		    product_id VARCHAR(36) NOT NULL,
		    quantity INT NOT NULL,
		    price BIGINT NOT NULL,
		    total BIGINT NOT NULL,
		    FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE RESTRICT,
		    FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT);
	`)
	if err != nil {
		panic(err)
	}

	return nil
}
