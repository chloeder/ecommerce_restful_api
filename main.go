package main

import (
	"database/sql"
	"log"
	"os"

	"commerce-project/migration"
	"commerce-project/routes"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database already connected")

	// Run migration
	err = migration.Migration(db)
	if err != nil {
		log.Fatal(err)
	}

	// Run routes
	err = routes.Routes(db)
	if err != nil {
		log.Fatal(err)
	}
}
