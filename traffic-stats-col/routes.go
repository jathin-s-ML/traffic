package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type RequestLog struct {
	Method      string `json:"method"`
	URL         string `json:"url"`
	StatusCode  int    `json:"status_code"`
	RequestSize int    `json:"request_size"`
}

type TrafficStats struct {
	TotalRequests   int     `json:"total_requests"`
	MostUsedMethod  string  `json:"most_used_method"`
	MostAccessedURL string  `json:"most_accessed_url"`
	AvgRequestSize  float64 `json:"avg_request_size"`
}

// ✅ Handles incoming data and stores it in the database
func CollectDataHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Failed to read request body", zap.Error(err))
		sendJSONResponse(w, map[string]interface{}{"error": "Failed to read request body"}, http.StatusInternalServerError)
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
		logger.Error("Database insertion failed", zap.Error(err))
		sendJSONResponse(w, map[string]interface{}{"error": "Database error"}, http.StatusInternalServerError)
		return
	}

	logger.Info("",
    zap.String("method", logEntry.Method),
    zap.String("url", logEntry.URL),
    zap.Int("status_code", logEntry.StatusCode),
    zap.Int("request_size", logEntry.RequestSize),
)


	sendJSONResponse(w, map[string]interface{}{"message": "Data received"}, http.StatusOK)
}

// ✅ Pagination Support for Getting Logs
func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationParams(r)

	method := r.URL.Query().Get("method")
	url := r.URL.Query().Get("url")
	status := r.URL.Query().Get("status")
	byteSize := r.URL.Query().Get("byte_size")

	logs, totalLogs, err := GetPaginatedTrafficLogs(method, url, status, byteSize, page, limit)
	if err != nil {
		logger.Error("Failed to retrieve logs", zap.Error(err))
		sendJSONResponse(w, map[string]interface{}{"error": "Failed to retrieve logs"}, http.StatusInternalServerError)
		return
	}

	logger.Info("Logs retrieved successfully",
		zap.Int("total_logs", totalLogs),
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	sendJSONResponse(w, map[string]interface{}{
		"total_logs": totalLogs,
		"page":       page,
		"limit":      limit,
		"logs":       logs,
	}, http.StatusOK)
}

// ✅ Handles logs retrieval by method
func GetLogsByMethodHandler(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Query().Get("method")
	if method == "" {
		logger.Warn("Missing query parameter", zap.String("parameter", "method"))
		sendJSONResponse(w, map[string]interface{}{"error": "Method query parameter is required"}, http.StatusBadRequest)
		return
	}

	logs, err := GetTrafficLogsByMethod(method)
	if err != nil {
		logger.Error("Failed to retrieve logs by method", zap.String("method", method), zap.Error(err))
		sendJSONResponse(w, map[string]interface{}{"error": "Failed to retrieve logs"}, http.StatusInternalServerError)
		return
	}

	logger.Info("Logs retrieved by method",
		zap.String("method", method),
		zap.Int("log_count", len(logs)),
	)

	sendJSONResponse(w, map[string]interface{}{"logs": logs}, http.StatusOK)
}

// ✅ Get Traffic Statistics
func GetTrafficStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := GetTrafficStats()
	if err != nil {
		logger.Error("Failed to retrieve traffic statistics", zap.Error(err))
		sendJSONResponse(w, map[string]interface{}{"error": "Failed to retrieve statistics"}, http.StatusInternalServerError)
		return
	}

	logger.Info("Traffic statistics retrieved successfully",
		zap.Int("total_requests", stats.TotalRequests),
		zap.String("most_used_method", stats.MostUsedMethod),
		zap.String("most_accessed_url", stats.MostAccessedURL),
		zap.Float64("avg_request_size", stats.AvgRequestSize),
	)

	sendJSONResponse(w, stats, http.StatusOK)
}

// ✅ Get Pagination Parameters
func getPaginationParams(r *http.Request) (int, int) {
	page := 1
	limit := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		} else {
			logger.Warn("Invalid page parameter, defaulting to 1", zap.String("provided_page", p))
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		} else {
			logger.Warn("Invalid limit parameter, defaulting to 10", zap.String("provided_limit", l))
		}
	}

	logger.Info("Pagination parameters parsed", zap.Int("page", page), zap.Int("limit", limit))
	return page, limit
}

// ✅ Send JSON Response (Helper Function)
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}
