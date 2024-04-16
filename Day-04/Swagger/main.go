package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
			Thanks: "Thank you!",
		})
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Error: message})
}

func main() {
	http.HandleFunc("/buy_candy", buyCandyHandler)
	log.Println("Starting server on :3333")
	log.Fatal(http.ListenAndServe(":3333", nil))
}