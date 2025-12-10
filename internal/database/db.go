package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectDB() error {
	dbUrl := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return err
	}

	DB = conn
	return nil
}
