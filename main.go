package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Response struct for returning JSON
type Response struct {
	CurrentTime string `json:"current_time"`
}

var db *sql.DB

func main() {
	// Set up database connection
	var err error
	dsn := "username:password@tcp(127.0.0.1:3306)/time_db" // Replace with your MySQL credentials and database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Verify the database connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	http.HandleFunc("/current-time", currentTimeHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Get current Toronto time
	torontoLocation, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Unable to load timezone", http.StatusInternalServerError)
		log.Printf("Error loading timezone: %v", err)
		return
	}

	currentTime := time.Now().In(torontoLocation)

	// Log current time to the database
	if err := logTimeToDatabase(currentTime); err != nil {
		http.Error(w, "Unable to log time to database", http.StatusInternalServerError)
		log.Printf("Error logging time to database: %v", err)
		return
	}

	// Prepare JSON response
	response := Response{CurrentTime: currentTime.Format("2006-01-02 15:04:05")}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func logTimeToDatabase(currentTime time.Time) error {
	query := "INSERT INTO time_log (timestamp) VALUES (?)"
	_, err := db.Exec(query, currentTime)
	return err
}
