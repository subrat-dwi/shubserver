package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() *pgxpool.Pool {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSql")
	return db
}
