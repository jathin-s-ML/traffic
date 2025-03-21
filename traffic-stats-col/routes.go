package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// RequestLog represents the data to be stored
type RequestLog struct {
	Method      string `json:"method"`
	URL         string `json:"url"`
	StatusCode  int    `json:"status_code"`
	RequestSize int    `json:"request_size"`
}

// TrafficStats represents summarized statistics of request logs
type TrafficStats struct {
	TotalRequests  int     `json:"total_requests"`
	MostUsedMethod string  `json:"most_used_method"`
	MostAccessedURL string `json:"most_accessed_url"`
	AvgRequestSize float64 `json:"avg_request_size"`
}

// CollectDataHandler handles API request logging
func CollectDataHandler(w http.ResponseWriter, r *http.Request) {
	// Open the log file in append mode
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer logFile.Close()

	// Create a logger that writes to the file
	logger := log.New(logFile, "", log.LstdFlags)

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Prepare log data
	logEntry := RequestLog{
		Method:      r.Method,
		URL:         r.URL.Path,
		StatusCode:  http.StatusOK, // Default 200 OK
		RequestSize: len(body),
	}

	// Store in PostgreSQL
	err = InsertTrafficLog(logEntry.Method, logEntry.URL, logEntry.StatusCode, logEntry.RequestSize)
	if err != nil {
		logger.Println("❌ Error inserting data:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Convert log entry to JSON format
	responseJSON, _ := json.MarshalIndent(logEntry, "", "  ")

	// Log data to logs.txt instead of printing to terminal
	logger.Println("✅ Stored Data:\n", string(responseJSON))

	// Send success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received"))
}

// GetLogsHandler handles retrieving stored logs from the database
func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	logs, err := GetTrafficLogs()
	if err != nil {
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}

	// Convert logs to JSON format
	responseJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format logs", http.StatusInternalServerError)
		return
	}

	// Send logs as JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// GetTrafficStatsHandler handles retrieving statistics of stored logs
func GetTrafficStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := GetTrafficStats()
	if err != nil {
		http.Error(w, "Failed to retrieve statistics", http.StatusInternalServerError)
		return
	}

	// Convert stats to JSON format
	responseJSON, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format statistics", http.StatusInternalServerError)
		return
	}

	// Send statistics as JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
