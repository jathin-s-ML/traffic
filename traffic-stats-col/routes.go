package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type RequestLog struct {
	Method      string `json:"method"`
	URL         string `json:"url"`
	StatusCode  int    `json:"status_code"`
	RequestSize int    `json:"request_size"`
}

type TrafficStats struct {
	TotalRequests  int     `json:"total_requests"`
	MostUsedMethod string  `json:"most_used_method"`
	MostAccessedURL string `json:"most_accessed_url"`
	AvgRequestSize float64 `json:"avg_request_size"`
}

func CollectDataHandler(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	logEntry := RequestLog{
		Method:      r.Method,
		URL:         r.URL.Path,
		StatusCode:  http.StatusOK, 
		RequestSize: len(body),
	}

	err = InsertTrafficLog(logEntry.Method, logEntry.URL, logEntry.StatusCode, logEntry.RequestSize)
	if err != nil {
		logger.Println("Error inserting data:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	responseJSON, _ := json.MarshalIndent(logEntry, "", "  ")
	logger.Println("Stored Data:\n", string(responseJSON))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received"))
}

func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	logs, err := GetTrafficLogs()
	if err != nil {
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func GetTrafficStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := GetTrafficStats()
	if err != nil {
		http.Error(w, "Failed to retrieve statistics", http.StatusInternalServerError)
		return
	}
	responseJSON, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format statistics", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
func TruncateLogsHandler(w http.ResponseWriter, r *http.Request) {
	err := TruncateTrafficLogs()
	if err != nil {
		http.Error(w, "Failed to truncate table", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Traffic logs table truncated successfully!"))
}
