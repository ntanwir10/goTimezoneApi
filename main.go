package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TimeResponse struct {
	CurrentTime string `json:"current_time"`
	Timezone    string `json:"timezone"`
}

var db *sql.DB

func main() {
	// Initialize database connection
	var err error
	db, err = sql.Open("mysql", "root:admin@tcp(localhost:3306)/timedb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Set up HTTP route with CORS middleware
	http.HandleFunc("/time", corsMiddleware(getTorontoTime))
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getTorontoTime(w http.ResponseWriter, r *http.Request) {
	// Load Toronto's timezone
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get current time in Toronto
	torontoTime := time.Now().In(location)

	// Log time to database
	_, err = db.Exec("INSERT INTO time_logs (request_time) VALUES (?)", torontoTime)
	if err != nil {
		log.Printf("Error logging to database: %v", err)
	}

	// Prepare response
	response := TimeResponse{
		CurrentTime: torontoTime.Format(time.RFC3339),
		Timezone:    location.String(),
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
