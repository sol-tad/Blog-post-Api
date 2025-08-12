package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Load the JWT secret key from environment variables
var JWT_ACCESS_TOKEN_SECRET = os.Getenv("JWT_ACCESS_TOKEN_SECRET")

// AuthMiddleware verifies the JWT access token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header value
		authHeader := c.GetHeader("Authorization")
		
		// Check if header is missing or does not start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		// Extract the token string by trimming the "Bearer " prefix and any whitespace
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		tokenStr = strings.TrimSpace(tokenStr)

		// Print token for debugging (remove in production)
		fmt.Println("Token--------------------:", tokenStr)

		// Parse the token using the secret key and validate signing method
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Verify that the signing method is HMAC (expected)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key bytes for signature verification
			return []byte(JWT_ACCESS_TOKEN_SECRET), nil
		})

		// Extract claims as a MapClaims type
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// If claims are invalid or token parsing failed, reject the request
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		// Store the claims information into the Gin context for downstream handlers
		c.Set("role", claims["role"])
		c.Set("id", claims["user_id"])
		c.Set("username", claims["username"])
		
		// Allow the request to proceed to the next middleware or handler
		c.Next()
	}
}

// AdminOnly middleware restricts access to users with "admin" role
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the "role" from Gin context (set by AuthMiddleware)
		role, exists := c.Get("role")

		// If role is missing or not "admin", reject with Forbidden status
		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		// Role is admin, proceed to the next handler
		c.Next()
	}
}
