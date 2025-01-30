package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "user=postgres password=mysecretpassword dbname=postgres host=172.18.0.2 port=5432 sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database successfully")
}
