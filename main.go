package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TimeResponse struct {
	CurrentTime  string `json:"current_time"`  // ISO format (RFC3339)
	ReadableTime string `json:"readable_time"` // Human readable format
	ReadableDate string `json:"readable_date"` // Human readable date
	Timezone     string `json:"timezone"`
}

var db *sql.DB

func main() {
	// Get database connection details from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "admin")
	dbName := getEnv("DB_NAME", "timedb")

	// Create database connection string
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)

	// Initialize database connection
	var err error
	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Wait for database to be ready
	for i := 0; i < 30; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Waiting for database... %v", err)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	// Set up HTTP route with both middlewares
	http.HandleFunc("/time", errorHandler(corsMiddleware(getTorontoTime)))
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getTorontoTime(w http.ResponseWriter, r *http.Request) {
	// Load Toronto's timezone
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		log.Printf("Timezone error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get current time in Toronto
	torontoTime := time.Now().In(location)

	// Log time to database
	_, err = db.Exec("INSERT INTO time_logs (request_time) VALUES (?)", torontoTime)
	if err != nil {
		log.Printf("Database error: %v", err)
		// Continue processing but log the error
	}

	// Prepare response with multiple time formats
	response := TimeResponse{
		CurrentTime:  torontoTime.Format(time.RFC3339),
		ReadableTime: torontoTime.Format("3:04 PM"),           // e.g., "2:30 PM"
		ReadableDate: torontoTime.Format("Monday, January 2"), // e.g., "Monday, March 14"
		Timezone:     location.String(),
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Response encoding error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			return
		}

		// Call the next handler
		next(w, r)
	}
}

// Add a custom error handler middleware
func errorHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		handler(w, r)
	}
}
