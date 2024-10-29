package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
)

func init() {
	const connectionUrl = "postgres://postgres:password@localhost:5432/go_comerce?sslmode=disable"
	conn, err := sql.Open("pgx", connectionUrl)
	if err != nil {
		fmt.Printf("We have an error %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Printf("Terjadi kesalahan: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Koneksi berhasil")
}

func main() {
	fmt.Println("Hello, world!")
}
