package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Global DB instance
var db *sql.DB

// ✅ Initialize Database Connection
func InitDB(cfg *Config) {
	var err error

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode,
	)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal("Database is not reachable", zap.Error(err))
	}

	logger.Info("Connected to PostgreSQL successfully")
}

// ✅ Insert a new traffic log
func InsertTrafficLog(method, url string, statusCode, requestSize int) error {
	query := `INSERT INTO request_logs (method, url, status_code, request_size) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, method, url, statusCode, requestSize)
	if err != nil {
		logger.Error("Failed to insert data",
			zap.String("method", method),
			zap.String("url", url),
			zap.Int("status_code", statusCode),
			zap.Int("request_size", requestSize),
			zap.Error(err),
		)
		return fmt.Errorf("failed to insert data: %v", err)
	}

	// logger.Info("Data inserted into database successfully",
	// 	zap.String("method", method),
	// 	zap.String("url", url),
	// 	zap.Int("status_code", statusCode),
	// 	zap.Int("request_size", requestSize),
	// )

	return nil
}

// ✅ Retrieve traffic logs by HTTP method
func GetTrafficLogsByMethod(method string) ([]RequestLog, error) {
	query := `SELECT method, url, status_code, request_size FROM request_logs WHERE method = $1`
	rows, err := db.Query(query, method)
	if err != nil {
		logger.Error("Failed to retrieve data", zap.String("method", method), zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve data: %v", err)
	}
	defer rows.Close()

	var logs []RequestLog
	for rows.Next() {
		var logEntry RequestLog
		if err := rows.Scan(&logEntry.Method, &logEntry.URL, &logEntry.StatusCode, &logEntry.RequestSize); err != nil {
			logger.Error("Failed to scan row", zap.Error(err))
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		logs = append(logs, logEntry)
	}

	logger.Info("Retrieved logs by method", zap.String("method", method), zap.Int("log_count", len(logs)))
	return logs, nil
}

// ✅ Retrieve paginated traffic logs
func GetPaginatedTrafficLogs(method, url, status, byteSize string, page, limit int) ([]RequestLog, int, error) {
	query := `SELECT method, url, status_code, request_size FROM request_logs WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM request_logs WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if method != "" {
		query += fmt.Sprintf(" AND method = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND method = $%d", argIndex)
		args = append(args, method)
		argIndex++
	}
	if url != "" {
		query += fmt.Sprintf(" AND url = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND url = $%d", argIndex)
		args = append(args, url)
		argIndex++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status_code = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status_code = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if byteSize != "" {
		query += fmt.Sprintf(" AND request_size = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND request_size = $%d", argIndex)
		args = append(args, byteSize)
		argIndex++
	}

	offset := (page - 1) * limit
	query += fmt.Sprintf(" ORDER BY method LIMIT %d OFFSET %d", limit, offset)

	var totalLogs int
	err := db.QueryRow(countQuery, args...).Scan(&totalLogs)
	if err != nil {
		logger.Error("Failed to get total count", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get total count: %v", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		logger.Error("Failed to retrieve paginated logs", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to retrieve data: %v", err)
	}
	defer rows.Close()

	var logs []RequestLog
	for rows.Next() {
		var logEntry RequestLog
		if err := rows.Scan(&logEntry.Method, &logEntry.URL, &logEntry.StatusCode, &logEntry.RequestSize); err != nil {
			logger.Error("Failed to scan row", zap.Error(err))
			return nil, 0, fmt.Errorf("failed to scan row: %v", err)
		}
		logs = append(logs, logEntry)
	}

	logger.Info("Retrieved paginated logs", zap.Int("page", page), zap.Int("limit", limit), zap.Int("total_logs", totalLogs))
	return logs, totalLogs, nil
}

// ✅ Retrieve traffic statistics
func GetTrafficStats() (TrafficStats, error) {
	var stats TrafficStats

	query := `
		SELECT 
			COUNT(*) AS total_requests,
			(SELECT method FROM request_logs GROUP BY method ORDER BY COUNT(*) DESC LIMIT 1) AS most_used_method,
			(SELECT url FROM request_logs GROUP BY url ORDER BY COUNT(*) DESC LIMIT 1) AS most_accessed_url,
			COALESCE(AVG(request_size), 0) AS avg_request_size
		FROM request_logs;
	`

	err := db.QueryRow(query).Scan(&stats.TotalRequests, &stats.MostUsedMethod, &stats.MostAccessedURL, &stats.AvgRequestSize)
	if err != nil {
		logger.Error("Failed to retrieve traffic stats", zap.Error(err))
		return stats, fmt.Errorf("failed to retrieve stats: %v", err)
	}

	logger.Info("Retrieved traffic stats",
		zap.Int("total_requests", stats.TotalRequests),
		zap.String("most_used_method", stats.MostUsedMethod),
		zap.String("most_accessed_url", stats.MostAccessedURL),
		zap.Float64("avg_request_size", stats.AvgRequestSize),
	)

	return stats, nil
}

// ✅ Truncate the traffic logs table
func TruncateTrafficLogs() error {
	query := `TRUNCATE TABLE request_logs`
	_, err := db.Exec(query)
	if err != nil {
		logger.Error("Failed to truncate table", zap.Error(err))
		return fmt.Errorf("failed to truncate table: %v", err)
	}

	logger.Warn("Traffic logs table truncated successfully")
	return nil
}
