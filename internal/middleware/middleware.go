package middleware

import (
	"fmt"
	"net/http"
	"strings"

	token "github.com/aborgas90/expense-tracker-api/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        fmt.Println("AUTH HEADER :: ",authHeader)
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := token.ValidateAccessToken(tokenStr)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        // simpan userID ke context
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
