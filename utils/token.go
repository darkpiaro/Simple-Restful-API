package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// getJWTSecret returns the JWT secret from environment variables
func getJWTSecret() []byte {
	// Load .env file if it exists
	godotenv.Load()

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback to default secret (not recommended for production)
		secret = "your-secret-key-change-this-in-production"
	}

	return []byte(secret)
}

// Claims structure for JWT
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken creates a new JWT token for the user
func GenerateToken(userID int, username string) (string, error) {
	// Create claims with user data and expiration time (24 hours)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "simple-restful-api",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractTokenFromHeader extracts Bearer token from Authorization header
func ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
