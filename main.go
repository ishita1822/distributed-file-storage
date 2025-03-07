package main

import (
	"fmt"
	"log"

	"distributed-file-storage/api"
	"distributed-file-storage/db"
)

func main() {
	fmt.Println("Starting File Storage Backend...")

	err := db.InitDB() // Initialize DB before anything else
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	api.StartServer()
}
