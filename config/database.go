package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_NAME"),
    )

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("DB connection error:", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal("DB ping error:", err)
    }

    DB = db
    fmt.Println("DB Connected")
}
