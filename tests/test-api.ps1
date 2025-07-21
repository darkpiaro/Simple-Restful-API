# PowerShell script to test the REST API

Write-Host "Testing REST API..." -ForegroundColor Green

# Test 1: Create a new user
Write-Host "`n1. Creating a new user..." -ForegroundColor Yellow
$body = @{
    username = "darkpiaro"
    password = "password123"
    fullname = "Dark Piaro"
    email = "darkpiaro@example.com"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/users" -Method POST -Headers @{"Content-Type"="application/json"} -Body $body
    Write-Host "✅ User created successfully!" -ForegroundColor Green
    $response | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Failed to create user: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2: Login to get JWT token
Write-Host "`n2. Logging in to get JWT token..." -ForegroundColor Yellow
$loginBody = @{
    username = "darkpiaro"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -Headers @{"Content-Type"="application/json"} -Body $loginBody
    Write-Host "✅ Login successful!" -ForegroundColor Green
    $token = $loginResponse.token
    Write-Host "JWT Token: $token" -ForegroundColor Cyan
} catch {
    Write-Host "❌ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    return
}

# Test 3: Get all users (protected route)
Write-Host "`n3. Getting all users (protected route)..." -ForegroundColor Yellow
try {
    $usersResponse = Invoke-RestMethod -Uri "http://localhost:8080/users" -Method GET -Headers @{"Authorization"="Bearer $token"}
    Write-Host "✅ Retrieved users successfully!" -ForegroundColor Green
    $usersResponse | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Failed to get users: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: Get specific user by ID
Write-Host "`n4. Getting user by ID..." -ForegroundColor Yellow
try {
    $userResponse = Invoke-RestMethod -Uri "http://localhost:8080/users/1" -Method GET -Headers @{"Authorization"="Bearer $token"}
    Write-Host "✅ Retrieved user successfully!" -ForegroundColor Green
    $userResponse | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Failed to get user: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🎉 API testing completed!" -ForegroundColor Green
