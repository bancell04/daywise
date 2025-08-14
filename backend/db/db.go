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
	// _, er := Pool.Exec(context.Background(), `
	//     DROP TABLE IF EXISTS tasks;
	//     DROP TABLE IF EXISTS categories;
	// `)
	// if er != nil {
	// 	log.Fatalf("Failed to drop tables: %v", er)
	// }

	_, errr := Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			color CHAR(7)
		)
    `)

	if errr != nil {
		log.Fatalf("Failed to create tasks table: %v", errr)
	}

	_, err := Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS tasks (
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
    		category_id INT REFERENCES categories(id) ON DELETE SET NULL,
            start TIMESTAMPTZ,
            "end" TIMESTAMPTZ
        )
    `)

	if err != nil {
		log.Fatalf("Failed to create tasks table: %v", err)
	}
}
