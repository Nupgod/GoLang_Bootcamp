package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"path"
)

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

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	// Command-line flags
	candyType := flag.String("k", "", "Candy type")
	candyCount := flag.Int("c", 0, "Candy count")
	money := flag.Int("m", 0, "Money")
	flag.Parse()

	if !isFlagPassed("k") || !isFlagPassed("c") || !isFlagPassed("m") {
		log.Fatalln("Wrong arguments")
	}

	// Load the client certificate and key
	certFilePath := path.Join("minica-key-pair", "client-cert", "cert.pem")
	keyFilePath := path.Join("minica-key-pair", "client-cert", "key.pem")
	clientCert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		log.Fatalf("Failed to load client certificate and key: %v", err)
	}

	// Load the CA certificate
	caCertPath := path.Join("minica-key-pair", "minica.pem")
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS for the client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Create the request body
	order := Order{
		Money:      *money,
		CandyType:  *candyType,
		CandyCount: *candyCount,
	}
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to marshal order: %v", err)
	}

	// Send the POST request
	resp, err := client.Post("https://server-cert:3333/buy_candy", "application/json", strings.NewReader(string(orderJSON)))
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Print the response
	if response.Error != "" {
		fmt.Println(response.Error)
	} else {
		fmt.Printf("%s Your change is %s\n", response.Thanks, response.Change)
	}
}
