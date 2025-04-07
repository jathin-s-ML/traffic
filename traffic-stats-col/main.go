package main

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize logger
	InitLogger()
	defer logger.Sync()

	// Clear logs file on startup
	if err := os.WriteFile("logs.txt", []byte{}, 0644); err != nil {
		logger.Fatal("Failed to clear logs file", zap.Error(err))
	}

	// Load configuration
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Validate port
	if cfg.Server.Port == "" {
		logger.Fatal("Port must be specified in config.yaml")
	}
	serverAddress := ":" + cfg.Server.Port

	// Initialize database
	InitDB(cfg)

	// Define routes
	http.HandleFunc("/collect", CollectDataHandler)
	http.HandleFunc("/logs", GetLogsHandler)
	http.HandleFunc("/stats", GetTrafficStatsHandler)
	http.HandleFunc("/logs/method", GetLogsByMethodHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to Traffic Stats Collector!",
		})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr: serverAddress,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Traffic Stats Collector is running", zap.String("address", serverAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server crashed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
