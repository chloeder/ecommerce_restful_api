package models

import (
	"database/sql"
	"errors"
)

type Product struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	IsDeleted *bool  `json:"is_deleted,omitempty"`
}

func SelectAllProducts(db *sql.DB) ([]Product, error) {
	if db == nil {
		return nil, errors.New("No Database Connected")
	}

	query := `SELECT id, name, price FROM products WHERE is_deleted = FALSE;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func SelectProductById(db *sql.DB, id string) (Product, error) {
	if db == nil {
		return Product{}, errors.New("No Database Connected")
	}

	query := `SELECT id, name, price FROM products WHERE id = $1`
	row := db.QueryRow(query, id)

	var product Product

	err := row.Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func InsertProduct(db *sql.DB, product Product) error {
	if db == nil {
		return errors.New("No Database Connected")
	}

	query := `INSERT INTO products (id, name, price) VALUES ($1, $2, $3)`

	_, err := db.Exec(query, product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProduct(db *sql.DB, product Product) error {
	if db == nil {
		return errors.New("No Database Connected")
	}

	query := `UPDATE products SET name = $1, price = $2 WHERE id = $3`

	_, err := db.Exec(query, product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func SoftDeletedProduct(db *sql.DB, product Product) error {
	if db == nil {
		return errors.New("No Database Connected")
	}

	query := `UPDATE products SET is_deleted = true WHERE id = $1`

	_, err := db.Exec(query, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func HardDeletedProduct(db *sql.DB, product Product) error {
	if db == nil {
		return errors.New("No Database Connected")
	}

	query := `DELETE FROM products WHERE id = $1`

	_, err := db.Exec(query, product.ID)
	if err != nil {
		return err
	}

	return nil
}
