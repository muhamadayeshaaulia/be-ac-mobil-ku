package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type AuthMiddleware struct {
	AuthClient *auth.Client
}

func InitAuthMiddleware() *AuthMiddleware {
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	if serviceAccountPath == "" {
		serviceAccountPath = "firebase-service-account.json"
	}

	ctx := context.Background()
	var app *firebase.App
	var err error

	if _, statErr := os.Stat(serviceAccountPath); statErr == nil {
		opt := option.WithCredentialsFile(serviceAccountPath)
		app, err = firebase.NewApp(ctx, nil, opt)
		if err != nil {
			log.Printf("error initializing Firebase app: %v", err)
		}
	} else {
		log.Printf("Firebase service account file not found at %s. Running without Firebase verification (Dev mode only).", serviceAccountPath)
	}

	var authClient *auth.Client
	if app != nil {
		authClient, err = app.Auth(ctx)
		if err != nil {
			log.Printf("error getting Firebase Auth client: %v", err)
		}
	}

	return &AuthMiddleware{AuthClient: authClient}
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

		idToken := parts[1]

		// Development Bypass Option
		if os.Getenv("APP_ENV") == "development" || am.AuthClient == nil {
			// Mock tokens for easy testing
			if strings.HasPrefix(idToken, "dev-token-") {
				role := "pelanggan"
				if strings.Contains(idToken, "pengelola") {
					role = "pengelola_bengkel"
				}
				uid := strings.TrimPrefix(idToken, "dev-token-")
				
				c.Set("user_uid", uid)
				c.Set("user_email", uid+"@example.com")
				c.Set("user_name", "Developer User " + uid)
				c.Set("user_role", role)
				c.Next()
				return
			}
		}

		if am.AuthClient == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Firebase authentication client not initialized"})
			c.Abort()
			return
		}

		// Verify ID Token via Firebase Admin SDK
		token, err := am.AuthClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			log.Printf("Firebase Auth Token verification failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired authorization token"})
			c.Abort()
			return
		}

		// Inject user context
		c.Set("user_uid", token.UID)
		
		email := ""
		if val, ok := token.Claims["email"]; ok {
			email = val.(string)
		}
		c.Set("user_email", email)

		name := ""
		if val, ok := token.Claims["name"]; ok {
			name = val.(string)
		}
		c.Set("user_name", name)

		c.Next()
	}
}
