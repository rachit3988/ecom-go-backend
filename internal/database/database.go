package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

var DB *pgxpool.Pool

func ConnectDB() error {
	if DB != nil {
		return nil
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return fmt.Errorf("DATABASE_URL missing")
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return err
	}

	DB = pool
	return nil
}
