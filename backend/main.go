package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Node struct {
	ID       string   `json:"id" gorm:"primaryKey"`
	Name     string   `json:"name"`
	IP       string   `json:"ip"`
	Position Position `json:"position" gorm:"embedded"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Server struct {
	db *gorm.DB
}

// sendJSONError sends a JSON error response
func sendJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// validateIP checks if the IP address is valid and returns a cleaned version
func validateIP(ipStr string) (string, error) {
	// Remove port if present
	if strings.Contains(ipStr, ":") {
		ipStr = strings.Split(ipStr, ":")[0]
	}

	// Parse IP address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", errors.New("invalid IP address format")
	}

	// Check if it's an IPv4 address
	if ip.To4() == nil {
		return "", errors.New("only IPv4 addresses are supported")
	}

	// Check if it's a private or loopback address
	if ip.IsPrivate() || ip.IsLoopback() {
		return ipStr, nil
	}

	// Check if it's a valid public IP
	if !ip.IsGlobalUnicast() {
		return "", errors.New("IP address must be a valid public or private address")
	}

	return ipStr, nil
}

// calculatePositionFromIP converts an IP address to a 3D position
func calculatePositionFromIP(ipStr string) Position {
	// Parse IP address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return Position{X: 0, Y: 0, Z: 0}
	}

	// Convert IP to 4 bytes
	ipBytes := ip.To4()
	if ipBytes == nil {
		return Position{X: 0, Y: 0, Z: 0}
	}

	// Normalize values to range [-5, 5]
	normalize := func(b byte) float64 {
		return (float64(b)/255.0)*10 - 5
	}

	// Use first three bytes for X, Y, Z coordinates
	// Last byte can be used for additional variation if needed
	return Position{
		X: normalize(ipBytes[0]),
		Y: normalize(ipBytes[1]),
		Z: normalize(ipBytes[2]),
	}
}

// calculateDistance calculates the Euclidean distance between two positions
func calculateDistance(p1, p2 Position) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// adjustPosition adjusts the position to maintain minimum distance from other nodes
func (s *Server) adjustPosition(newPos Position, existingNodes []Node) Position {
	const minDistance = 1.0 // Minimum distance between nodes
	adjustedPos := newPos

	for _, node := range existingNodes {
		dist := calculateDistance(adjustedPos, node.Position)
		if dist < minDistance {
			// Calculate direction vector
			dx := adjustedPos.X - node.Position.X
			dy := adjustedPos.Y - node.Position.Y
			dz := adjustedPos.Z - node.Position.Z

			// Normalize and scale
			length := math.Sqrt(dx*dx + dy*dy + dz*dz)
			if length > 0 {
				scale := (minDistance - dist) / length
				adjustedPos.X += dx * scale
				adjustedPos.Y += dy * scale
				adjustedPos.Z += dz * scale
			}
		}
	}

	// Ensure position stays within bounds [-5, 5]
	clamp := func(v float64) float64 {
		return math.Max(-5, math.Min(5, v))
	}

	return Position{
		X: clamp(adjustedPos.X),
		Y: clamp(adjustedPos.Y),
		Z: clamp(adjustedPos.Z),
	}
}

func (s *Server) handleGetNodes(w http.ResponseWriter, r *http.Request) {
	var nodes []Node
	if err := s.db.Find(&nodes).Error; err != nil {
		sendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func (s *Server) handleCreateNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var input struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendJSONError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate IP address
	validIP, err := validateIP(input.IP)
	if err != nil {
		log.Printf("Invalid IP address: %s - %v", input.IP, err)
		sendJSONError(w, http.StatusBadRequest, "Invalid IP address: "+err.Error())
		return
	}

	// Check if node with name already exists
	var existingNodeWithName Node
	if err := s.db.Where("name = ?", input.Name).First(&existingNodeWithName).Error; err == nil {
		log.Printf("Node name already taken: %s", input.Name)
		sendJSONError(w, http.StatusConflict, "Node name already taken")
		return
	}

	// Check if node with IP already exists
	var existingNodeWithIP Node
	if err := s.db.Where("ip = ?", validIP).First(&existingNodeWithIP).Error; err == nil {
		log.Printf("Node already exists for IP: %s", validIP)
		sendJSONError(w, http.StatusConflict, "IP node already exists")
		return
	}

	// Get all existing nodes for position adjustment
	var existingNodes []Node
	s.db.Find(&existingNodes)

	// Calculate initial position from IP
	initialPos := calculatePositionFromIP(validIP)

	// Adjust position to avoid overlaps
	finalPos := s.adjustPosition(initialPos, existingNodes)

	// Create new node
	node := Node{
		ID:       generateID(),
		Name:     input.Name,
		IP:       validIP,
		Position: finalPos,
	}

	if err := s.db.Create(&node).Error; err != nil {
		log.Printf("Error creating node: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Failed to create node: "+err.Error())
		return
	}

	log.Printf("Successfully created node: %s for IP: %s at position: %+v", node.ID, validIP, finalPos)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Get database configuration from environment variables
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "password")
	dbHost := getEnvOrDefault("DB_HOST", "mariadb")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "ipnetwork")

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&Node{})

	// Create server
	server := &Server{db: db}

	// Create router
	mux := http.NewServeMux()
	
	// API v1 routes
	v1 := http.NewServeMux()
	v1.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.handleGetNodes(w, r)
		case http.MethodPost:
			server.handleCreateNode(w, r)
		default:
			sendJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	v1.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			sendJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		server.handleGetIP(w, r)
	})

	// Mount v1 routes under /v1 prefix
	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(mux)))
}

// getEnvOrDefault returns the value of the environment variable or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// handleGetIP handles the request to get the client's public IP
func (s *Server) handleGetIP(w http.ResponseWriter, r *http.Request) {
	// Get the X-Real-IP header first (set by nginx/proxy)
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// Try X-Forwarded-For header
		ip = r.Header.Get("X-Forwarded-For")
		if ip != "" {
			// If X-Forwarded-For contains multiple IPs, take the first one
			ip = strings.Split(ip, ",")[0]
		}
	}
	
	// If no proxy headers, use RemoteAddr
	if ip == "" {
		ip = r.RemoteAddr
	}

	// Remove port if present
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}

	// Validate the IP
	validIP, err := validateIP(ip)
	if err != nil {
		log.Printf("Invalid IP from request: %s - %v", ip, err)
		sendJSONError(w, http.StatusBadRequest, "Invalid IP address: "+err.Error())
		return
	}

	// Return the IP
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"ip": validIP,
	})
}
