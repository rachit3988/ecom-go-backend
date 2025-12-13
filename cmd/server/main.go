package main

import (
	"ecom-go-backend/api"
	"ecom-go-backend/internal/database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer database.DB.Close()

	api.SetupRoutes()
	fmt.Println("Server running on http://localhost:8080")
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
