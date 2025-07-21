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

// User represents a user in the system
type User struct {
	UserID    int    `json:"userid" example:"1"`
	Username  string `json:"username" example:"johndoe"`
	Password  string `json:"password,omitempty"` // omitempty เพื่อไม่ส่ง password ใน response
	FullName  string `json:"fullname" example:"John Doe"`
	Email     string `json:"email" example:"john@example.com"`
	CreatedBy *int   `json:"created_by,omitempty"`
	CreatedOn string `json:"created_on,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
	UpdatedOn string `json:"updated_on,omitempty"`
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

	// Drop and recreate table with new structure
	dropTableQuery := `IF EXISTS (SELECT * FROM sysobjects WHERE name='users' AND xtype='U') DROP TABLE users`
	_, err = db.Exec(dropTableQuery)
	if err != nil {
		log.Printf("Warning: Could not drop existing table: %v", err)
	}

	// Create users table with new structure
	createTableQuery := `
	CREATE TABLE users (
		userid INT IDENTITY(1,1) PRIMARY KEY,
		username NVARCHAR(50) UNIQUE NOT NULL,
		password NVARCHAR(255) NOT NULL,
		fullname NVARCHAR(100) NOT NULL,
		email NVARCHAR(100) NULL,
		created_by INT NULL,
		created_on DATETIME DEFAULT GETDATE(),
		updated_by INT NULL,
		updated_on DATETIME NULL
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

	query := "INSERT INTO users (username, password, fullname, email) OUTPUT INSERTED.userid VALUES (@username, @password, @fullname, @email)"
	var newID int
	err = db.QueryRow(query,
		sql.Named("username", u.Username),
		sql.Named("password", string(hashedPassword)),
		sql.Named("fullname", u.FullName),
		sql.Named("email", u.Email)).Scan(&newID)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	u.UserID = newID
	u.Password = "" // Clear password from struct
	return nil
}

// GetAll retrieves all users
func GetAllUsers() ([]User, error) {
	query := "SELECT userid, username, fullname, email, created_on FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Username, &user.FullName, &user.Email, &user.CreatedOn)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// GetByID retrieves a user by ID
func GetUserByID(id int) (*User, error) {
	query := "SELECT userid, username, fullname, email, created_on FROM users WHERE userid = @id"
	row := db.QueryRow(query, sql.Named("id", id))

	var user User
	err := row.Scan(&user.UserID, &user.Username, &user.FullName, &user.Email, &user.CreatedOn)
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
	query := "SELECT userid, username, password, fullname, email FROM users WHERE username = @username"
	row := db.QueryRow(query, sql.Named("username", username))

	var user User
	err := row.Scan(&user.UserID, &user.Username, &user.Password, &user.FullName, &user.Email)
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
		query = "UPDATE users SET username = @username, password = @password, fullname = @fullname, email = @email, updated_on = GETDATE() WHERE userid = @id"
		_, err = db.Exec(query,
			sql.Named("username", u.Username),
			sql.Named("password", string(hashedPassword)),
			sql.Named("fullname", u.FullName),
			sql.Named("email", u.Email),
			sql.Named("id", u.UserID))
	} else {
		// Update without password
		query = "UPDATE users SET username = @username, fullname = @fullname, email = @email, updated_on = GETDATE() WHERE userid = @id"
		_, err = db.Exec(query,
			sql.Named("username", u.Username),
			sql.Named("fullname", u.FullName),
			sql.Named("email", u.Email),
			sql.Named("id", u.UserID))
	}

	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	u.Password = "" // Clear password from struct
	return nil
}

// Delete deletes a user
func DeleteUser(id int) error {
	query := "DELETE FROM users WHERE userid = @id"
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
