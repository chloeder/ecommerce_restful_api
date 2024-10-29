package db

import (
	"database/sql"
	"fmt"
	"os"
)

func init() {
	conn, err := sql.Open("pgx", "postgresql://postgres:password@localhost:5432/ecommerce_api?sslmode=disable")
	if err != nil {
		fmt.Printf("Server Error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Printf("Server Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database already connected")

	// Create table products
	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS products (
    		id VARCHAR(36) PRIMARY KEY,
    		name VARCHAR(255) NOT NULL,
    		price BIGINT NOT NULL,
    		is_deleted BOOLEAN DEFAULT FALSE
    )`)
	if err != nil {
		fmt.Printf("Server Error: %v\n", err)
		os.Exit(1)
	}
}
