package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"os"
)

const (
	DbDriver   = "DB_DRIVER"
	DbUserName = "DB_USER"
	DbPassword = "DB_PASSWORD"
	DbHostName = "DB_HOST"
	DbName     = "DB_NAME"
)

// InitDB initializes the database connection and returns the *sql.DB instance
func InitDB() (*sql.DB, error) {

	dbDriver, _ := os.LookupEnv(DbDriver)
	dbUserName, _ := os.LookupEnv(DbUserName)
	dbPassword, _ := os.LookupEnv(DbPassword)
	dbHostName, _ := os.LookupEnv(DbHostName)
	dbName, _ := os.LookupEnv(DbName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUserName, dbPassword, dbHostName, dbName)

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Connected to the database successfully")
	return db, nil
}
