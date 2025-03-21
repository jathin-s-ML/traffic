package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

// APIRequest interface
type APIRequest interface {
	SendRequest(url string) error
}

// Define request types
type GetRequest struct{}
type PostRequest struct{}
type PutRequest struct{}
type DeleteRequest struct{}

// Implement SendRequest for each request type
func (g GetRequest) SendRequest(url string) error {
	return sendHTTPRequest("GET", url, nil)
}

func (p PostRequest) SendRequest(url string) error {
	payload := RandomData()
	return sendHTTPRequest("POST", url, payload)
}

func (p PutRequest) SendRequest(url string) error {
	payload := RandomData()
	return sendHTTPRequest("PUT", url, payload)
}

func (d DeleteRequest) SendRequest(url string) error {
	return sendHTTPRequest("DELETE", url, nil)
}

// Function to send HTTP requests and log details
func sendHTTPRequest(method, url string, body []byte) error {
	client := &http.Client{}
	var req *http.Request
	var err error

	// Capture request body size
	bodySize := len(body)

	if method == "POST" || method == "PUT" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Capture response status
	statusCode := resp.StatusCode

	// Prepare log entry
	logEntry := fmt.Sprintf("[Request] Method: %s, URL: %s, Body Size: %d bytes\n[Response] Status: %d\n",
		method, url, bodySize, statusCode)

	// Write log entry to file
	err = writeLog(logEntry)
	if err != nil {
		fmt.Println("❌ Error writing to log file:", err)
	}

	return nil
}

// Function to write logs to log.txt
func writeLog(entry string) error {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(entry)
	return err
}

// Function to get a random API request type
func GetRandomRequest() APIRequest {
	requests := []APIRequest{GetRequest{}, PostRequest{}, PutRequest{}, DeleteRequest{}}
	return requests[rand.Intn(len(requests))]
}
