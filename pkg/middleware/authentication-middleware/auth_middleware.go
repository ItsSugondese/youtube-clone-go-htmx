package authentication_middleware

import (
	"youtube-clone/pkg/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function that checks for the presence of a valid token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Here you can add your authentication logic, e.g., checking for a token in the request header
		token := c.GetHeader("Authorization")

		if token == "" || !verifyTokenService(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// If the token is valid, proceed to the next handler
		c.Next()
	}
}

func verifyTokenService(token string) bool {
	publicKey := os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	if publicKey == "" {
		panic("public key not found in environment")
	}

	return utils.ValidateToken(token[7:], publicKey)
}
