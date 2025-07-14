package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	connStr := os.Getenv("POSTGRES_DSN")
	if connStr == "" {
		return fmt.Errorf("POSTGRES_DSN not set in environment")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}
