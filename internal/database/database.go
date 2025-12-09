package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() error {
	connStr := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return err
	}

	DB = pool
	return nil
}
