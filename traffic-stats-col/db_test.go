package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Load environment variables from .env file
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// Get database connection from environment variables
func getDBConnection() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	return sql.Open("postgres", dsn)
}

func TestInsertData(t *testing.T) {
	db, err := getDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	query := `INSERT INTO request_logs (method, url, status_code, request_size) 
			  VALUES ('GET', '/test', 200, 123) RETURNING id;`

	var id int
	err = db.QueryRow(query).Scan(&id)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	t.Logf("Inserted row with ID: %d", id)
}

func TestRetrieveLogs(t *testing.T) {
	db, err := getDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, method, url, status_code, request_size FROM request_logs LIMIT 1")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var id, statusCode, requestSize int
	var method, url string

	if rows.Next() {
		err := rows.Scan(&id, &method, &url, &statusCode, &requestSize)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		t.Logf("Retrieved log: ID=%d, Method=%s, URL=%s, Status=%d, Size=%d", id, method, url, statusCode, requestSize)
	} else {
		t.Log("No logs found")
	}
}

func TestGetStats(t *testing.T) {
	db, err := getDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	var totalRequests int
	var mostUsedMethod, mostAccessedURL string
	var avgRequestSize float64

	query := `SELECT COUNT(*), 
	                 (SELECT method FROM request_logs GROUP BY method ORDER BY COUNT(*) DESC LIMIT 1),
	                 (SELECT url FROM request_logs GROUP BY url ORDER BY COUNT(*) DESC LIMIT 1),
	                 AVG(request_size) 
	          FROM request_logs;`

	err = db.QueryRow(query).Scan(&totalRequests, &mostUsedMethod, &mostAccessedURL, &avgRequestSize)
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	stats := map[string]interface{}{
		"total_requests":   totalRequests,
		"most_used_method": mostUsedMethod,
		"most_accessed_url": mostAccessedURL,
		"avg_request_size": avgRequestSize,
	}

	statsJSON, _ := json.Marshal(stats)
	t.Logf("Stats: %s", statsJSON)
}

func TestTruncateTable(t *testing.T) {
	db, err := getDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("TRUNCATE TABLE request_logs RESTART IDENTITY;")
	if err != nil {
		t.Fatalf("Truncate failed: %v", err)
	}
	t.Log("Table truncated successfully")
}