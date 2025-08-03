package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	dsn := "postgres://daywiseuser:daywisepassword@localhost:5432/daywisedb"
	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Connected to DB!")
}

func Close() {
	Pool.Close()
}

func Setup() {
	_, err := Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS tasks (
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
            category TEXT,
            start TIMESTAMPTZ,
            "end" TIMESTAMPTZ
        )
    `)

	if err != nil {
		log.Fatalf("Failed to create tasks table: %v", err)
	}
}
