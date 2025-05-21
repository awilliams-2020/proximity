package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Node struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"position"`
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Test IP addresses (IPv4 only)
	testIPs := []string{
		// AWS US East
		"3.5.140.0",
		"3.5.141.0",
		"3.5.142.0",
		// AWS US West
		"13.56.0.0",
		"13.57.0.0",
		"13.58.0.0",
		// Google Cloud
		"34.95.0.0",
		"34.96.0.0",
		"34.97.0.0",
		// Azure
		"20.0.0.0",
		"20.1.0.0",
		"20.2.0.0",
		// Digital Ocean
		"143.198.0.0",
		"143.198.1.0",
		"143.198.2.0",
		// Public DNS
		"8.8.8.8",        // Google DNS
		"1.1.1.1",        // Cloudflare DNS
		"9.9.9.9",        // Quad9 DNS
		"208.67.222.222", // OpenDNS
	}

	// Create HTTP client
	client := &http.Client{}

	// Create nodes with different IPs
	for i, ip := range testIPs {
		// Create request body
		body := map[string]string{
			"name": fmt.Sprintf("Node-%d", i+1),
		}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			log.Fatalf("Error marshaling request body: %v", err)
		}

		// Create request
		req, err := http.NewRequest("POST", "http://localhost:8080/nodes", bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", ip)

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request for IP %s: %v", ip, err)
			continue
		}

		// Read response
		var node Node
		if err := json.NewDecoder(resp.Body).Decode(&node); err != nil {
			log.Printf("Error decoding response for IP %s: %v", ip, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// Check response
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error creating node for IP %s: status code %d", ip, resp.StatusCode)
			continue
		}

		log.Printf("Successfully created node: %s with IP: %s at position: %+v",
			node.Name, node.IP, node.Position)

		// Add a small delay between requests
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Test completed")
}
