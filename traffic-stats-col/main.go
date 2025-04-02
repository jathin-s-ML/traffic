package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
)

// var logger *zap.Logger

func main() {
	// Initialize Uber Zap logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}
	defer logger.Sync() // Flush logs before exit

	// Clear logs file on startup
	err = os.WriteFile("logs.txt", []byte{}, 0644)
	if err != nil {
		logger.Fatal("Failed to clear logs file", zap.Error(err))
	}

	// Load configuration
	cfg, err := ReadConfig()
	if err != nil {
		logger.Fatal("Failed to read config", zap.Error(err))
	}

	// Initialize database
	InitDB(cfg)

	// Define routes
	http.HandleFunc("/collect", CollectDataHandler)        // Collect request data
	http.HandleFunc("/logs", GetLogsHandler)              // Get logs with pagination
	http.HandleFunc("/stats", GetTrafficStatsHandler)     // Get traffic statistics
	http.HandleFunc("/logs/method", GetLogsByMethodHandler) // Get logs filtered by method

	// Start the server
	logger.Info("Traffic Stats Collector is running",
		zap.String("port", cfg.Server.Port),
	)
	logger.Fatal("Server crashed", zap.Error(http.ListenAndServe(":"+cfg.Server.Port, nil)))
}
