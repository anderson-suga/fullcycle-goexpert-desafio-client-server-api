package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Currency represents the specific JSON structure returned by the external API.
type Currency struct {
	USDBRL struct {
		Code       string `json:"code"`
		CodeIn     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

// BidOutput represents the simplified JSON response sent to the client.
type BidOutput struct {
	Bid string `json:"bid"`
}

func main() {
	// Initialize SQLite database connection
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// Register handler and inject database dependency
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		GetExchangeRateHandler(w, r, db)
	})

	// Start HTTP server
	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

// GetExchangeRateHandler handles the request for dollar quotation.
func GetExchangeRateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// 1. Fetch data from external API with a 200ms timeout
	currency, err := getExternalCurrency()
	if err != nil {
		log.Printf("Error fetching API: %v", err)
		http.Error(w, "Error fetching exchange rate", http.StatusRequestTimeout)
		return
	}

	// 2. Persist data to SQLite with a 10ms timeout
	err = saveCurrency(db, currency.USDBRL.Bid)
	if err != nil {
		log.Printf("Error saving to database: %v", err)
		http.Error(w, "Error saving exchange rate to database", http.StatusInternalServerError)
		return
	}

	// 3. Return result to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(BidOutput{Bid: currency.USDBRL.Bid})
}

// getExternalCurrency performs the HTTP request to the external API with context control.
func getExternalCurrency() (*Currency, error) {
	// Create a context with a 200ms timeout
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Create the HTTP request with the context
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	// Perform the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into Currency struct
	var currency Currency
	err = json.Unmarshal(body, &currency)
	if err != nil {
		return nil, err
	}

	return &currency, nil
}

// saveCurrency inserts the bid value into the database with context control.
func saveCurrency(db *sql.DB, bid string) error {
	// Create a context with a 10ms timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Prepare and execute the insert statement
	stmt, err := db.Prepare("INSERT INTO cotacoes(bid) VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement with context
	_, err = stmt.ExecContext(ctx, bid)
	if err != nil {
		return err
	}

	return nil
}