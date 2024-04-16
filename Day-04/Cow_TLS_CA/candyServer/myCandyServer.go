package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

// CandyPrices represents the prices of different candy types.
var CandyPrices = map[string]int{
	"CE": 10,
	"AA": 15,
	"NT": 17,
	"DE": 21,
	"YR": 23,
}

// Order represents a candy order.
type Order struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

// Response represents the server's response.
type Response struct {
	Change string `json:"change,omitempty"` // Use pointer to int to omit when zero
	Thanks string `json:"thanks,omitempty"`
	Error  string `json:"error,omitempty"`
}

func buyCandyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the order
	if order.CandyCount < 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Candy count cannot be negative")
		return
	}

	if order.Money < 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Money value cannot be negative")
		return
	}

	price, ok := CandyPrices[order.CandyType]
	if !ok {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid candy type")
		return
	}

	totalPrice := price * order.CandyCount
	if totalPrice > order.Money {
		neededMoney := totalPrice - order.Money
		writeErrorResponse(w, http.StatusPaymentRequired, fmt.Sprintf("You need %d more money!", neededMoney))
		return
	}

	// If everything is okay, respond with 201 and the change
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	change := order.Money - totalPrice
	json.NewEncoder(w).Encode(Response{
		Change: strconv.Itoa(change),
		Thanks: ask_cow(),
	})
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Error: message})
}

func main() {
	// Load the server certificate and key
	certFilePath := path.Join("cow-key-pair", "server-cert", "cert.pem")
	if _, err := os.Stat(certFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Failed to finde file certificate: %v", err)
	}
	keyFilePath := path.Join("cow-key-pair", "server-cert", "key.pem")
	if _, err := os.Stat(keyFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Failed to finde file key: %v", err)
	}
	serverCert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		log.Fatalf("Failed to load server certificate and key: %v", err)
	}

	// Load the CA certificate
	caCertPath := path.Join("cow-key-pair", "minica.pem")
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS for the server
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	// Create a server with the TLS config
	server := &http.Server{
		Addr:      ":3333",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/buy_candy", buyCandyHandler)
	log.Println("Starting server on :3333")

	// Start the server with TLS
	log.Fatal(server.ListenAndServeTLS("", ""))
}
