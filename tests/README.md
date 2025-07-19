# API Test Suite

This folder contains comprehensive test scripts for the Simple RESTful API.

## 📁 Test Files

### **PowerShell Test Scripts:**
- `test-api.ps1` - Basic API functionality tests
- `test-crud.ps1` - Complete CRUD operation tests
- `test-errors.ps1` - Security and error handling tests

### **Documentation:**
- `TEST_RESULTS.md` - Comprehensive test results and status report

## 🚀 How to Run Tests

### **Prerequisites:**
1. Ensure the API server is running (`go run main.go` from project root)
2. Open PowerShell in the project root directory

### **Running Individual Tests:**

```powershell
# Basic functionality test
.\tests\test-api.ps1

# Complete CRUD operations test
.\tests\test-crud.ps1

# Security and error handling test
.\tests\test-errors.ps1
```

### **Running All Tests:**

```powershell
# Run all tests in sequence
.\tests\test-api.ps1; .\tests\test-crud.ps1; .\tests\test-errors.ps1
```

## 📊 Test Coverage

### **✅ Basic API Tests (`test-api.ps1`):**
- User creation
- User authentication (login)
- JWT token generation
- Get all users (protected route)
- Get single user by ID (protected route)

### **✅ CRUD Tests (`test-crud.ps1`):**
- User update (PUT)
- Password change verification
- User deletion (DELETE)
- Verification of deletion
- Final state check

### **✅ Error Handling Tests (`test-errors.ps1`):**
- Unauthorized access attempts
- Invalid login credentials
- Duplicate username creation
- Non-existent user access
- Invalid JSON format handling

## 🔧 Test Environment

- **API Server**: `http://localhost:8080`
- **Database**: `smart-village-system.duckdns.org`
- **Authentication**: JWT Bearer tokens
- **Test Data**: Creates and manages test users automatically

## 📈 Expected Results

All tests should pass with:
- ✅ Green checkmarks for successful operations
- ✅ Proper error codes for security tests
- ✅ Consistent API response times
- ✅ Proper JWT token functionality

## 🛠 Troubleshooting

If tests fail:
1. Ensure the server is running on port 8080
2. Check database connectivity
3. Verify `.env` configuration
4. Check PowerShell execution policy: `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser`

## 📝 Adding New Tests

To add new test scenarios:
1. Follow the existing PowerShell script patterns
2. Use proper error handling with try/catch blocks
3. Include colored output for test results
4. Update this README with new test descriptions
