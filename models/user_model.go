package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // omitempty เพื่อไม่ส่ง password ใน response
	FullName string `json:"full_name"`
}

var db *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using system environment variables")
	}

	// Get database configuration from environment variables
	dbServer := os.Getenv("DB_SERVER")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Set default values if environment variables are not set
	if dbServer == "" {
		dbServer = "localhost"
	}
	if dbUser == "" {
		dbUser = "sa"
	}
	if dbPassword == "" {
		dbPassword = "YourPassword123"
	}
	if dbPort == "" {
		dbPort = "1433"
	}
	if dbName == "" {
		dbName = "TestDB"
	}

	// Build connection string from environment variables
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s",
		dbServer, dbUser, dbPassword, dbPort, dbName)

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	// Create users table if not exists
	createTableQuery := `
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='users' AND xtype='U')
	CREATE TABLE users (
		id INT IDENTITY(1,1) PRIMARY KEY,
		username NVARCHAR(50) UNIQUE NOT NULL,
		password NVARCHAR(255) NOT NULL,
		full_name NVARCHAR(100) NOT NULL
	)`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// Create creates a new user
func (u *User) Create() error {
	// Hash password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	query := "INSERT INTO users (username, password, full_name) OUTPUT INSERTED.id VALUES (@username, @password, @fullname)"
	var newID int
	err = db.QueryRow(query,
		sql.Named("username", u.Username),
		sql.Named("password", string(hashedPassword)),
		sql.Named("fullname", u.FullName)).Scan(&newID)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	u.ID = newID
	u.Password = "" // Clear password from struct
	return nil
}

// GetAll retrieves all users
func GetAllUsers() ([]User, error) {
	query := "SELECT id, username, full_name FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.FullName)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// GetByID retrieves a user by ID
func GetUserByID(id int) (*User, error) {
	query := "SELECT id, username, full_name FROM users WHERE id = @id"
	row := db.QueryRow(query, sql.Named("id", id))

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.FullName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %v", err)
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username (for login)
func GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, username, password, full_name FROM users WHERE username = @username"
	row := db.QueryRow(query, sql.Named("username", username))

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.FullName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %v", err)
	}

	return &user, nil
}

// Update updates a user
func (u *User) Update() error {
	var query string
	var err error

	// If password is provided, hash it and update
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing password: %v", err)
		}
		query = "UPDATE users SET username = @username, password = @password, full_name = @fullname WHERE id = @id"
		_, err = db.Exec(query,
			sql.Named("username", u.Username),
			sql.Named("password", string(hashedPassword)),
			sql.Named("fullname", u.FullName),
			sql.Named("id", u.ID))
	} else {
		// Update without password
		query = "UPDATE users SET username = @username, full_name = @fullname WHERE id = @id"
		_, err = db.Exec(query,
			sql.Named("username", u.Username),
			sql.Named("fullname", u.FullName),
			sql.Named("id", u.ID))
	}

	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	u.Password = "" // Clear password from struct
	return nil
}

// Delete deletes a user
func DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = @id"
	result, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ValidatePassword checks if the provided password matches the hashed password
func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
