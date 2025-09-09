package main

import (
	"context"
	"ecom-go-backend/api"
	"ecom-go-backend/internal/database"
	"ecom-go-backend/internal/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var db *pgx.Conn

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	var err error
	db, err = database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close(context.Background())

	// test ping
	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Successfully pinged database!")

	handlers.Db = db

	api.SetupRoutes()
	fmt.Println("Server running on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
