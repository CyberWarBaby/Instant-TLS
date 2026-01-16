package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instanttls/api/internal/config"
	"github.com/instanttls/api/internal/models"
	"github.com/jmoiron/sqlx"
)

// PATAuth validates Personal Access Token from Authorization header
func PATAuth(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		tokenHash := hashToken(token)

		var tokenRecord models.Token
		err := db.Get(&tokenRecord, `
			SELECT id, user_id, name, prefix, token_hash, last_used_at, created_at
			FROM tokens WHERE token_hash = $1
		`, tokenHash)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Update last_used_at
		go func() {
			db.Exec("UPDATE tokens SET last_used_at = $1 WHERE id = $2", time.Now(), tokenRecord.ID)
		}()

		// Get user
		var user models.User
		err = db.Get(&user, "SELECT * FROM users WHERE id = $1", tokenRecord.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("token", tokenRecord)
		c.Next()
	}
}

// SessionAuth validates session cookie for web dashboard
func SessionAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// For MVP, we use a simple JWT-like token in cookie or header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Try cookie
			authHeader, _ = c.Cookie("auth_token")
		}

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Remove Bearer prefix if present
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Decode the simple session token (user_id:email:plan encoded)
		parts := strings.Split(token, ":")
		if len(parts) != 3 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		user := models.User{
			Email: parts[1],
			Plan:  models.Plan(parts[2]),
		}

		// Parse UUID
		if err := user.ID.UnmarshalText([]byte(parts[0])); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
