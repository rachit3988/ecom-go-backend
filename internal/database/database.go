package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn, error) {
	// Use environment variables for credentials
	connStr := os.Getenv("DATABASE_URL")
	return pgx.Connect(context.Background(), connStr)
}
