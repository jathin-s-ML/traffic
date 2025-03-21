-- Create database
CREATE DATABASE traffic_data_20_mar;

-- Connect to the database
\c traffic_data_20_mar;

-- Create the request_logs table
CREATE TABLE IF NOT EXISTS request_logs (
    id SERIAL PRIMARY KEY,
    method VARCHAR(10) NOT NULL,
    url TEXT NOT NULL,
    request_size INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
