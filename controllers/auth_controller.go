package controllers

import (
	"net/http"
	"simple-restful-api/models"
	"simple-restful-api/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"darkpiaro"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token   string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User    models.User `json:"user"`
	Message string      `json:"message" example:"Login successful"`
}

// Login handles user authentication
// @Summary User Login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 401 {object} map[string]interface{} "Invalid username or password"
// @Failure 500 {object} map[string]interface{} "Failed to generate token"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginReq LoginRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get user from database
	user, err := models.GetUserByUsername(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Validate password
	err = user.ValidatePassword(loginReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Clear password from user object before sending response
	user.Password = ""

	// Send success response with token
	response := LoginResponse{
		Token:   token,
		User:    *user,
		Message: "Login successful",
	}

	c.JSON(http.StatusOK, response)
}
