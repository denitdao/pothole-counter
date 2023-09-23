package util

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

// DB is a package level variable that will be available for
// sharing between different packages in the application.
var DB *sql.DB

func InitDB() {
	driver := GetProperty("db.driver")
	dataSourceName := GetProperty("db.datasource")

	var err error
	DB, err = sql.Open(driver, dataSourceName)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Set connection pool parameters
	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
}
