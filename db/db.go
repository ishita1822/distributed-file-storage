package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := "host=localhost port=5432 user=admin password=password dbname=file_storage sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
		return err
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database is not reachable:", err)
		return err
	}

	fmt.Println("Connected to PostgreSQL successfully!")
	return nil
}
