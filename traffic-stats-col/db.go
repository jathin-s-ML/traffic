package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(cfg *Config) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	fmt.Println("Connected to PostgreSQL successfully!")
}

func InsertTrafficLog(method, url string, statusCode, requestSize int) error {
	query := `INSERT INTO request_logs (method, url, status_code, request_size) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, method, url, statusCode, requestSize)
	if err != nil {
		return fmt.Errorf("failed to insert data: %v", err)
	}
	fmt.Println("Data inserted into database successfully!")
	return nil
}

func GetTrafficLogs() ([]RequestLog, error) {
	rows, err := db.Query("SELECT method, url, status_code, request_size FROM request_logs")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data: %v", err)
	}
	defer rows.Close()

	var logs []RequestLog
	for rows.Next() {
		var logEntry RequestLog
		if err := rows.Scan(&logEntry.Method, &logEntry.URL, &logEntry.StatusCode, &logEntry.RequestSize); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		logs = append(logs, logEntry)
	}

	return logs, nil
}

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
		return stats, fmt.Errorf("failed to retrieve stats: %v", err)
	}

	return stats, nil
}
func TruncateTrafficLogs() error {
	query := `TRUNCATE TABLE request_logs`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to truncate table: %v", err)
	}
	fmt.Println("Traffic logs table truncated successfully!")
	return nil
}
func GetTrafficLogsByMethod(method string) ([]RequestLog, error) {
	query := `SELECT method, url, status_code, request_size FROM request_logs WHERE method = $1`
	rows, err := db.Query(query, method)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data: %v", err)
	}
	defer rows.Close()

	var logs []RequestLog
	for rows.Next() {
		var logEntry RequestLog
		if err := rows.Scan(&logEntry.Method, &logEntry.URL, &logEntry.StatusCode, &logEntry.RequestSize); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		logs = append(logs, logEntry)
	}

	return logs, nil
}
func GetFilteredTrafficLogs(method, url, status, byteSize string) ([]RequestLog, error) {
	query := `SELECT method, url, status_code, request_size FROM request_logs WHERE 1=1`
	args := []interface{}{}
	argIndex := 1
	if method != "" {
		query += fmt.Sprintf(" AND method = $%d", argIndex)
		args = append(args, method)
		argIndex++
	}
	if url != "" {
		query += fmt.Sprintf(" AND url = $%d", argIndex)
		args = append(args, url)
		argIndex++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status_code = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if byteSize != "" {
		query += fmt.Sprintf(" AND request_size = $%d", argIndex)
		args = append(args, byteSize)
		argIndex++
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data: %v", err)
	}
	defer rows.Close()

	var logs []RequestLog
	for rows.Next() {
		var logEntry RequestLog
		if err := rows.Scan(&logEntry.Method, &logEntry.URL, &logEntry.StatusCode, &logEntry.RequestSize); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		logs = append(logs, logEntry)
	}

	return logs, nil
}
