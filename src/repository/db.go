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
	initFloorSchema()
	initMenuSchema()
	initOrderSchema()
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

func initFloorSchema() {
	query := `
    CREATE TABLE IF NOT EXISTS floors (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        description TEXT,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        "order" INT NOT NULL
    );
    `
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create floor schema: %v", err)
	}

	log.Println("Floor schema initialized successfully")
}

func initMenuSchema() {
	query := `
    CREATE TABLE IF NOT EXISTS menus (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        price FLOAT NOT NULL,
        recipe TEXT,
        categories TEXT[],
        image VARCHAR(255),
        description TEXT,
        availability BOOLEAN NOT NULL DEFAULT FALSE,
        estimated_time INT,
        ingredients TEXT[],
        created_by INT REFERENCES users(id),
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );
    `
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create menu schema: %v", err)
	}

	log.Println("Menu schema initialized successfully")
}

func initOrderSchema() {
	query := `
    CREATE TABLE IF NOT EXISTS orders (
        id SERIAL PRIMARY KEY,
        table_id INT REFERENCES tables(id),
        user_id INT REFERENCES users(id),
        status VARCHAR(20),
        notes TEXT,
        payment_method VARCHAR(20),
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        total_amount FLOAT,
        discount FLOAT,
        tax FLOAT
    );
    CREATE TABLE IF NOT EXISTS order_items (
        id SERIAL PRIMARY KEY,
        order_id INT REFERENCES orders(id),
        menu_id INT REFERENCES menus(id),
        quantity INT,
        price FLOAT,
        created_at TIMESTAMP NOT NULL
    );
    `
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create order schema: %v", err)
	}

	log.Println("Order schema initialized successfully")
}
