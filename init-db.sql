DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'traffic_data_20_mar') THEN
        CREATE DATABASE traffic_data_20_mar;
    END IF;
END $$;

-- Connect to the database
\c traffic_data_20_mar;

CREATE TABLE IF NOT EXISTS request_logs (
    id SERIAL PRIMARY KEY,
    method VARCHAR(10) NOT NULL,
    url TEXT NOT NULL,
    status_code INT NOT NULL,
    request_size INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
