package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_ACCESS_TOKEN_SECRET = os.Getenv("JWT_ACCESS_TOKEN_SECRET")

func AuthMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
      return
    }

    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    tokenStr = strings.TrimSpace(tokenStr) 

	fmt.Println("Token--------------------:", tokenStr)

    token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
      // Optional: Verify signing method
      if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
      }
      return []byte(JWT_ACCESS_TOKEN_SECRET), nil
    })


    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
      return
    }

    c.Set("role", claims["role"])
    c.Set("id", claims["user_id"])
    c.Set("username", claims["username"])
    c.Next()
  }
}


func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role,exists:=c.Get("role")
		if!exists||role!="admin"{
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		c.Next()

	}
}