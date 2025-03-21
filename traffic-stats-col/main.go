package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Clear the logs file on startup
	err := os.WriteFile("logs.txt", []byte{}, 0644)
	if err != nil {
		log.Fatal("❌ Failed to clear logs file:", err)
	}

	// Load configuration
	cfg, err := ReadConfig()
	if err != nil {
		log.Fatal("❌ Failed to read config:", err)
	}

	// Initialize database connection
	InitDB(cfg)

	// Register routes
	http.HandleFunc("/collect", CollectDataHandler)
	http.HandleFunc("/logs", GetLogsHandler) // API to fetch stored logs
	http.HandleFunc("/stats", GetTrafficStatsHandler) // New API to get statistics

	// Start HTTP server
	fmt.Println("🚀 Traffic Stats Collector is running on port", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, nil))
}
