package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	// Ambil DSN dari variabel lingkungan
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Jika variabel lingkungan tidak diatur, gunakan default
		dsn = "root@tcp(127.0.0.1:3306)/learn_db"
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected!")
}
