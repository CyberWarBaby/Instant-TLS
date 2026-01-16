package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/CyberWarBaby/Instant-TLS/apps/api/internal/config"
	"github.com/CyberWarBaby/Instant-TLS/apps/api/internal/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db     *sqlx.DB
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func New(db *sqlx.DB, cfg *config.Config, logger *zap.SugaredLogger) *Handler {
	return &Handler{db: db, cfg: cfg, logger: logger}
}

// Register creates a new user
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var count int
	h.db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", req.Email)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Errorf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Create user
	userID := uuid.New()
	_, err = h.db.Exec(`
		INSERT INTO users (id, email, password_hash, plan)
		VALUES ($1, $2, $3, $4)
	`, userID, req.Email, string(passwordHash), "free")

	if err != nil {
		h.logger.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create session token
	sessionToken := createSessionToken(userID, req.Email, "free")

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   sessionToken,
		"user": models.UserResponse{
			ID:        userID,
			Email:     req.Email,
			Plan:      models.PlanFree,
			CreatedAt: time.Now(),
		},
	})
}

// Login authenticates a user
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	var user models.User
	err := h.db.Get(&user, "SELECT * FROM users WHERE email = $1", req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create session token
	sessionToken := createSessionToken(user.ID, user.Email, string(user.Plan))

	c.JSON(http.StatusOK, gin.H{
		"token": sessionToken,
		"user": models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Plan:      user.Plan,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Me returns current user info (PAT auth)
func (h *Handler) Me(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	c.JSON(http.StatusOK, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Plan:      user.Plan,
		CreatedAt: user.CreatedAt,
	})
}

// GetUser returns current user info (session auth)
func (h *Handler) GetUser(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	// Get full user from DB
	var fullUser models.User
	err := h.db.Get(&fullUser, "SELECT * FROM users WHERE id = $1", user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		ID:        fullUser.ID,
		Email:     fullUser.Email,
		Plan:      fullUser.Plan,
		CreatedAt: fullUser.CreatedAt,
	})
}

// UpdatePlan updates user's subscription plan
func (h *Handler) UpdatePlan(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var req struct {
		Plan string `json:"plan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate plan
	validPlans := map[string]bool{"free": true, "pro": true, "team": true}
	if !validPlans[req.Plan] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan. Must be 'free', 'pro', or 'team'"})
		return
	}

	// Update plan in database
	_, err := h.db.Exec(`UPDATE users SET plan = $1 WHERE id = $2`, req.Plan, user.ID)
	if err != nil {
		h.logger.Errorf("Failed to update plan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plan"})
		return
	}

	// Get updated user
	var updatedUser models.User
	h.db.Get(&updatedUser, "SELECT * FROM users WHERE id = $1", user.ID)

	c.JSON(http.StatusOK, models.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Plan:      updatedUser.Plan,
		CreatedAt: updatedUser.CreatedAt,
	})
}

// License returns user's license info
func (h *Handler) License(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	limits := map[string]int{
		"max_wildcard_certs": 1,
	}

	if user.Plan == models.PlanPro || user.Plan == models.PlanTeam {
		limits["max_wildcard_certs"] = -1 // unlimited
	}

	c.JSON(http.StatusOK, models.LicenseResponse{
		Plan:   user.Plan,
		Limits: limits,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Plan:      user.Plan,
			CreatedAt: user.CreatedAt,
		},
	})
}

// ListTokens returns all tokens for the user
func (h *Handler) ListTokens(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var tokens []models.Token
	err := h.db.Select(&tokens, `
		SELECT id, user_id, name, prefix, token_hash, last_used_at, created_at
		FROM tokens WHERE user_id = $1 ORDER BY created_at DESC
	`, user.ID)

	if err != nil {
		h.logger.Errorf("Failed to list tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tokens"})
		return
	}

	response := make([]models.TokenResponse, len(tokens))
	for i, t := range tokens {
		response[i] = models.TokenResponse{
			ID:         t.ID,
			Name:       t.Name,
			Prefix:     t.Prefix,
			LastUsedAt: t.LastUsedAt,
			CreatedAt:  t.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateToken creates a new Personal Access Token
func (h *Handler) CreateToken(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		h.logger.Errorf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	token := "itls_" + hex.EncodeToString(tokenBytes)
	prefix := token[:12]
	tokenHash := hashToken(token)

	tokenID := uuid.New()
	_, err := h.db.Exec(`
		INSERT INTO tokens (id, user_id, name, prefix, token_hash)
		VALUES ($1, $2, $3, $4, $5)
	`, tokenID, user.ID, req.Name, prefix, tokenHash)

	if err != nil {
		h.logger.Errorf("Failed to create token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusCreated, models.TokenCreateResponse{
		Token: token, // Only shown once!
		Data: models.TokenResponse{
			ID:        tokenID,
			Name:      req.Name,
			Prefix:    prefix,
			CreatedAt: time.Now(),
		},
	})
}

// DeleteToken revokes a token
func (h *Handler) DeleteToken(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	tokenID := c.Param("id")

	result, err := h.db.Exec(`
		DELETE FROM tokens WHERE id = $1 AND user_id = $2
	`, tokenID, user.ID)

	if err != nil {
		h.logger.Errorf("Failed to delete token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token revoked successfully"})
}

// MachinePing registers or updates a machine
func (h *Handler) MachinePing(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var req struct {
		Hostname string `json:"hostname" binding:"required"`
		OS       string `json:"os" binding:"required"`
		Arch     string `json:"arch" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Upsert machine
	_, err := h.db.Exec(`
		INSERT INTO machines (id, user_id, hostname, os, arch, last_seen_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, hostname) DO UPDATE SET
			os = EXCLUDED.os,
			arch = EXCLUDED.arch,
			last_seen_at = EXCLUDED.last_seen_at
	`, uuid.New(), user.ID, req.Hostname, req.OS, req.Arch, time.Now())

	if err != nil {
		h.logger.Errorf("Failed to ping machine: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register machine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Machine registered successfully"})
}

func createSessionToken(userID uuid.UUID, email, plan string) string {
	return userID.String() + ":" + email + ":" + plan
}

func hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
