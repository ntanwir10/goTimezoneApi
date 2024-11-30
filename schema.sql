-- create the SQL schema for the database:

CREATE DATABASE timedb;

USE timedb;

CREATE TABLE time_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    request_time DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 