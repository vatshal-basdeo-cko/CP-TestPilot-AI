package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testpilot-ai/gateway/auth"
	"github.com/testpilot-ai/shared/logger"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	db *pgxpool.Pool
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{db: db}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Query user from database
	var userID uuid.UUID
	var passwordHash, role string

	requestID, _ := c.Get("request_id")
	requestIDStr, _ := requestID.(string)

	query := "SELECT id, password_hash, role FROM users WHERE username = $1"
	err := h.db.QueryRow(c.Request.Context(), query, req.Username).Scan(&userID, &passwordHash, &role)
	if err != nil {
		logger.WithRequestID(requestIDStr).Debug().
			Str("username", req.Username).
			Msg("Login failed: user not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, passwordHash) {
		logger.WithRequestID(requestIDStr).Debug().
			Str("username", req.Username).
			Str("user_id", userID.String()).
			Msg("Login failed: invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(userID, role)
	if err != nil {
		logger.WithRequestID(requestIDStr).Err(err).
			Str("username", req.Username).
			Str("user_id", userID.String()).
			Msg("Failed to generate token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	logger.WithRequestID(requestIDStr).Info().
		Str("username", req.Username).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("User logged in successfully")

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       userID,
			"username": req.Username,
			"role":     role,
		},
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	requestID, _ := c.Get("request_id")
	requestIDStr, _ := requestID.(string)

	// Password validation - minimum 8 characters
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
		return
	}

	// Role validation - public registration can only create "user" role
	// Admin role can only be assigned by existing admins via CreateUser endpoint
	if req.Role != "" && req.Role != "user" {
		logger.WithRequestID(requestIDStr).Warn().
			Str("username", req.Username).
			Str("attempted_role", req.Role).
			Msg("Attempted to self-assign non-user role during registration")
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot self-assign admin role. Use admin panel to create admin users."})
		return
	}
	req.Role = "user" // Always set to user for public registration

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user
	userID := uuid.New()
	query := `
		INSERT INTO users (id, username, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	_, err = h.db.Exec(c.Request.Context(), query, userID, req.Username, passwordHash, req.Role, now, now)
	if err != nil {
		logger.WithRequestID(requestIDStr).Debug().
			Str("username", req.Username).
			Msg("Registration failed: username already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(userID, req.Role)
	if err != nil {
		logger.WithRequestID(requestIDStr).Err(err).
			Str("username", req.Username).
			Str("user_id", userID.String()).
			Msg("Failed to generate token after registration")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User created but token generation failed"})
		return
	}

	logger.WithRequestID(requestIDStr).Info().
		Str("username", req.Username).
		Str("user_id", userID.String()).
		Str("role", req.Role).
		Msg("User registered successfully")

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":       userID,
			"username": req.Username,
			"role":     req.Role,
		},
	})
}

// Me returns current user info
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	role := c.MustGet("role").(string)

	var username string
	query := "SELECT username FROM users WHERE id = $1"
	h.db.QueryRow(c.Request.Context(), query, userID).Scan(&username)

	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
		"role":     role,
	})
}

// ListUsers returns all users (admin only)
func (h *AuthHandler) ListUsers(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	query := `SELECT id, username, role, created_at FROM users ORDER BY created_at DESC`
	rows, err := h.db.Query(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []gin.H
	for rows.Next() {
		var id uuid.UUID
		var username, userRole string
		var createdAt time.Time
		if err := rows.Scan(&id, &username, &userRole, &createdAt); err != nil {
			continue
		}
		users = append(users, gin.H{
			"id":         id,
			"username":   username,
			"role":       userRole,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// CreateUser creates a new user (admin only) - allows setting any role
func (h *AuthHandler) CreateUser(c *gin.Context) {
	// Check if requester is admin
	requesterRole := c.MustGet("role").(string)
	if requesterRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	requestID, _ := c.Get("request_id")
	requestIDStr, _ := requestID.(string)

	// Password validation
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
		return
	}

	// Default role if not specified
	if req.Role == "" {
		req.Role = "user"
	}

	// Validate role is either "user" or "admin"
	if req.Role != "user" && req.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'user' or 'admin'"})
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user
	userID := uuid.New()
	query := `
		INSERT INTO users (id, username, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	_, err = h.db.Exec(c.Request.Context(), query, userID, req.Username, passwordHash, req.Role, now, now)
	if err != nil {
		logger.WithRequestID(requestIDStr).Debug().
			Str("username", req.Username).
			Msg("User creation failed: username already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	logger.WithRequestID(requestIDStr).Info().
		Str("created_by", c.MustGet("user_id").(uuid.UUID).String()).
		Str("username", req.Username).
		Str("role", req.Role).
		Msg("User created by admin")

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":       userID,
			"username": req.Username,
			"role":     req.Role,
		},
	})
}

// DeleteUser deletes a user (admin only)
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	// Check if requester is admin
	requesterRole := c.MustGet("role").(string)
	if requesterRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Prevent self-deletion
	requesterID := c.MustGet("user_id").(uuid.UUID)
	if userID == requesterID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	requestID, _ := c.Get("request_id")
	requestIDStr, _ := requestID.(string)

	// Delete user
	query := `DELETE FROM users WHERE id = $1`
	result, err := h.db.Exec(c.Request.Context(), query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	logger.WithRequestID(requestIDStr).Info().
		Str("deleted_by", requesterID.String()).
		Str("deleted_user_id", userIDStr).
		Msg("User deleted by admin")

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"id":      userIDStr,
	})
}




