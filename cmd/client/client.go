package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// BidOutput matches the JSON structure returned by our server.
type BidOutput struct {
	Bid string `json:"bid"`
}

func main() {
	// Define a context with a 300ms timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Prepare the request with the context
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Could not create request: %v", err)
	}

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request (timeout likely): %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		// Read the error message from the response body
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Server returned error: %s", string(body))
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}

	// Unmarshal JSON into BidOutput struct
	var bid BidOutput
	err = json.Unmarshal(body, &bid)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Save to file
	err = saveToFile(bid.Bid)
	if err != nil {
		log.Fatalf("Error saving file: %v", err)
	}

	fmt.Println("Process finished successfully. Bid saved.")
}

// saveToFile writes the result to cotacao.txt in the format "Dólar: {value}".
func saveToFile(bid string) error {
	content := fmt.Sprintf("Dólar: %s", bid)
	
	// 0644 provides read/write permissions for the owner and read for others
	err := os.WriteFile("cotacao.txt", []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}