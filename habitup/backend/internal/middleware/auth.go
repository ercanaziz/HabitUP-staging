package middleware

import (
	"context"
	"net/http"
	"strings"

	"habitup/internal/auth"
	"habitup/internal/cache"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token gerekli"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		if cache.IsBlacklisted(context.Background(), token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token geçersiz (logout yapılmış)"})
			return
		}

		claims, err := auth.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Geçersiz token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
