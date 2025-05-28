package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"report-service/config"            // Project's config package
)

// InitDB initializes and returns a database connection
func InitDB(cfg *config.Config) (*sql.DB, error) {
	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Open a new database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
