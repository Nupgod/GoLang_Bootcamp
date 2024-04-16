package jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Define a secret key for JWT token generation (this should be securely stored and not exposed in code)
var jwtKey = []byte("my_secret_key")

// Claims defines the structure of the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// TokenResponse defines the structure of the response for the /api/get_token endpoint
type TokenResponse struct {
	Token string `json:"token"`
}

// generateTokenHandler generates a JWT token for a user
func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// For simplicity, we'll just generate a token for a hardcoded user
	username := "user" // This should be replaced with actual authentication logic

	// Define the expiration time of the token
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	json.NewEncoder(w).Encode(TokenResponse{Token: tokenString})
}

// middleware to authenticate JWT token
func AuthenticateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Split the header to get the Bearer token
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// Parse the JWT token
		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next(w, r)
	}
}
