// example_usage.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"your/path/to/httpdbg"
)

// LoginResponse represents the structure of the login response
type LoginResponse struct {
	Token string `json:"token"`
}

func main() {
	// Create a client with debug logging
	client := httpdbg.NewClient()

	// Prepare login request (replace with your actual login endpoint and credentials)
	loginURL := "https://your-cyberark-api.com/login"

	// You might need to adjust this based on your specific API requirements
	loginPayload := []byte(`{
		"username": "your_username",
		"password": "your_password"
	}`)

	// Perform login request
	resp, err := client.Post(loginURL, "application/json", bytes.NewBuffer(loginPayload))
	if err != nil {
		log.Fatalf("Login request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		log.Fatalf("Failed to parse login response: %v", err)
	}

	// Print the extracted token
	fmt.Println("Extracted Token:", loginResp.Token)

	// Use the token in a subsequent request
	// This is just an example - replace with your actual API endpoint
	req, err := http.NewRequest("GET", "https://your-cyberark-api.com/some-endpoint", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", loginResp.Token)

	// Send the request using the same debug client
	apiResp, err := client.Do(req)
	if err != nil {
		log.Fatalf("API request failed: %v", err)
	}
	defer apiResp.Body.Close()
}
