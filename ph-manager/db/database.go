package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"ph-manager/util"
	"time"
)

// DB is a connection pool to the database.
var DB *sql.DB

func InitDB() {
	driver := util.GetProperty("db.driver")
	dataSourceName := util.GetProperty("db.datasource")

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
