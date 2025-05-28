package main

import (
	"log"
	"report-service/config"
	"report-service/database"
	"report-service/report"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // MySQL driver, ensure it's imported for side effects
)

func main() {
	// 2.a. Load application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2.b. Initialize database connection
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	// 2.c. Create a new Gin router
	router := gin.Default()

	// 2.d. Register the report export handler
	// No specific route group mentioned, so registering directly.
	// The handler expects *sql.DB. The config (cfg) is not directly passed to ExportReportHandler
	// as per its current definition: report.ExportReportHandler(db *sql.DB)
	router.GET("/report/export", report.ExportReportHandler(db))

	// 2.e. Start the Gin server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
