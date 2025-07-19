# Environment Configuration Setup

## Changes Made

The project has been updated to use environment variables from a `.env` file for configuration. This makes the application more secure and flexible for different deployment environments.

### Files Modified:

1. **go.mod** - Added `github.com/joho/godotenv v1.5.1` dependency
2. **models/user_model.go** - Updated to read database configuration from environment variables
3. **utils/token.go** - Updated to read JWT secret from environment variables
4. **main.go** - Updated to read server port from environment variables

### Environment Variables Used:

```
# Database Configuration
DB_SERVER=smart-village-system.duckdns.org
DB_USER=sa
DB_PASSWORD=P@ssw0rd
DB_PORT=1433
DB_NAME=Go_Simple_DB

# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production

# Server Configuration
PORT=8080
```

### Key Features:

1. **Automatic .env loading** - The application automatically loads configuration from `.env` file
2. **Fallback defaults** - If environment variables are not set, the application uses sensible defaults
3. **Security** - Sensitive information like passwords and JWT secrets are kept in environment variables
4. **Flexibility** - Easy to change configuration for different environments (development, staging, production)

### Usage:

1. Copy `.env.example` to `.env`
2. Update the values in `.env` file according to your environment
3. Run the application with `go run main.go`

### Before Running:

Make sure to install dependencies:
```bash
go mod tidy
```

The application will now use the database configuration from your `.env` file:
- Server: `smart-village-system.duckdns.org`
- Database: `Go_Simple_DB`
- User: `sa`
- Password: `P@ssw0rd`
- Port: `1433`
