package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	JWTSecret []byte
}

func InitAuthMiddleware() *AuthMiddleware {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "ac_mobil_ku_secret_key_1234567890" // default fallback
	}
	return &AuthMiddleware{JWTSecret: []byte(secret)}
}

func (am *AuthMiddleware) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in format 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Development Bypass Option
		if os.Getenv("APP_ENV") == "development" {
			// Mock tokens for easy offline testing
			if strings.HasPrefix(tokenString, "dev-token-") {
				role := "pelanggan"
				if strings.Contains(tokenString, "pengelola") {
					role = "pengelola_bengkel"
				}
				uid := strings.TrimPrefix(tokenString, "dev-token-")
				
				c.Set("user_uid", uid)
				c.Set("user_email", uid+"@example.com")
				c.Set("user_name", "Developer User " + uid)
				c.Set("user_role", role)
				c.Next()
				return
			}
		}

		// Parse and validate custom JWT
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return am.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired JWT token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT claims structure"})
			c.Abort()
			return
		}

		// Extract claims into context
		uid, _ := claims["uid"].(string)
		if uid == "" {
			uid, _ = claims["sub"].(string) // fallback to standard subject claim
		}
		
		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)
		role, _ := claims["role"].(string)

		c.Set("user_uid", uid)
		c.Set("user_email", email)
		c.Set("user_name", name)
		c.Set("user_role", role)

		c.Next()
	}
}
