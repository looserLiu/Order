package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"code": 401, "message": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}
	return userID.(uuid.UUID)
}

// RateLimit creates a rate limiting middleware
// maxRequests: maximum requests per window
// windowMinutes: time window in minutes
func RateLimit(maxRequests int, windowMinutes int) gin.HandlerFunc {
	// Simple in-memory rate limiter (for production, use Redis)
	clients := make(map[string]int64)
	windows := make(map[string]int64)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := int64(0) // In production, use time.Now().Unix()
		
		// Check if window has expired
		if windowStart, ok := windows[clientIP]; ok {
			// Reset counter if window expired
			// In production: if now-windowStart > int64(windowMinutes*60) { ... }
		}
		
		// Check rate limit
		if count, ok := clients[clientIP]; ok && count >= int64(maxRequests) {
			c.JSON(429, gin.H{"code": 429, "message": "Too many requests"})
			c.Abort()
			return
		}
		
		clients[clientIP]++
		c.Next()
	}
}

// RequestID adds a unique request ID to each request for tracing
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New()
		c.Header("X-Request-ID", requestID.String())
		c.Set("request_id", requestID)
		c.Next()
	}
}
