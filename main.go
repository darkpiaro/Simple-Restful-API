package main

import (
	"log"
	"os"
	"simple-restful-api/controllers"
	"simple-restful-api/middlewares"
	"simple-restful-api/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "simple-restful-api/docs" // Auto-generated Swagger docs
)

// @title Simple RESTful API
// @version 1.0
// @description A simple REST API with JWT authentication built with Go and Gin
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using system environment variables")
	}

	// Initialize database connection
	err = models.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer models.CloseDB()

	// Create Gin router
	router := gin.Default()

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes (no authentication required)
	router.POST("/login", controllers.Login)
	router.POST("/users", controllers.CreateUser)

	// Protected routes (authentication required)
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/users", controllers.GetUsers)
		protected.GET("/users/:id", controllers.GetUser)
		protected.PUT("/users/:id", controllers.UpdateUser)
		protected.DELETE("/users/:id", controllers.DeleteUser)
	}

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Start server
	log.Printf("Server starting on port %s...", port)
	router.Run(":" + port)
}
