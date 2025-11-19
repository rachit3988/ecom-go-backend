package main

import (
	"context"
	"ecom-go-backend/api"
	"ecom-go-backend/internal/database"
	"ecom-go-backend/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var db *pgx.Conn

func init() {
	_ = godotenv.Load()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default for local dev
	}

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
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
