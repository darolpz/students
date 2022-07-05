package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/darolpz/students/internal/auth"
	"github.com/gin-gonic/gin"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrParseClaims   = errors.New("couldn't parse claims")
	ErrInvalidClaims = errors.New("claims are not valid")
)

func AuthMiddleware(authService auth.IAuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get the token from the header
		authorization := c.GetHeader("Authorization")

		// If there is no token, return error
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrTokenNotFound.Error()})
			return
		}

		// Split the token into the Bearer and the token
		authorization = strings.Split(authorization, "Bearer ")[1]

		// Check token
		token, err := authService.CheckToken(authorization)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Parse claims into JWTClaim struct
		claims, ok := token.Claims.(*auth.JWTClaim)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrParseClaims.Error()})
			return
		}

		// Check if claims are valid
		if err := claims.Valid(); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidClaims.Error()})
			return
		}

		// Set the user to the context
		c.Header("email", claims.Email)
		c.Header("role", claims.Role)
		c.Next()
	}
}
