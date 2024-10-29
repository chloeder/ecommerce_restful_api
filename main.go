package main

import (
	"commerce-project/routes"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	routes.Routes()
}
