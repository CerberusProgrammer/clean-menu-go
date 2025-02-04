package repository

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connection established successfully")

	initUserSchema()
	initTableSchema()
}

func initUserSchema() {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) NOT NULL,
        password VARCHAR(100) NOT NULL,
        name VARCHAR(50),
        last_name VARCHAR(50),
        email VARCHAR(100) UNIQUE NOT NULL,
        phone VARCHAR(20),
        role VARCHAR(20),
        image VARCHAR(255)
    );
    `
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create user schema: %v", err)
	}

	log.Println("User table schema initialized successfully")
}

func initTableSchema() {
	query := `
    CREATE TABLE IF NOT EXISTS tables (
        id SERIAL PRIMARY KEY,
        number VARCHAR(50) NOT NULL,
        name VARCHAR(50),
        capacity INT NOT NULL,
        shape VARCHAR(20),
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        status VARCHAR(20)
    );
    `
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table schema: %v", err)
	}

	log.Println("Table schema initialized successfully")
}
