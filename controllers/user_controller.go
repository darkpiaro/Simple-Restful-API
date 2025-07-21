package controllers

import (
	"net/http"
	"simple-restful-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
	FullName string `json:"full_name" binding:"required" example:"John Doe"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Username string `json:"username" example:"johndoe_updated"`
	Password string `json:"password" example:"newpassword123"`
	FullName string `json:"full_name" example:"John Doe Updated"`
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user account (public endpoint)
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User creation data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 500 {object} map[string]interface{} "Failed to create user"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var req CreateUserRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create user model
	user := models.User{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
	}

	// Save user to database
	err := user.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUsers retrieves all users
// @Summary Get all users
// @Description Retrieve a list of all users (protected endpoint)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of users"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve users"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve users",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// GetUser retrieves a single user by ID
// @Summary Get user by ID
// @Description Retrieve a specific user by their ID (protected endpoint)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User details"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	// Get user ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Get user from database
	user, err := models.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateUser updates an existing user
// @Summary Update user
// @Description Update an existing user's information (protected endpoint)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User update data"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Failed to update user"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	// Get user ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Check if user exists
	existingUser, err := models.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	// Bind JSON request body
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Update user fields if provided
	if req.Username != "" {
		existingUser.Username = req.Username
	}
	if req.FullName != "" {
		existingUser.FullName = req.FullName
	}
	if req.Password != "" {
		existingUser.Password = req.Password
	}

	// Save updated user
	err = existingUser.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    existingUser,
	})
}

// DeleteUser deletes a user by ID
// @Summary Delete user
// @Description Delete a user by their ID (protected endpoint)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	// Get user ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Delete user from database
	err = models.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
