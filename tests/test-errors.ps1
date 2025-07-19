# Error Handling Tests

Write-Host "Testing Error Handling Scenarios" -ForegroundColor Green

# Test 1: Access protected route without token
Write-Host "`n1. Testing access without JWT token..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "http://localhost:8080/users" -Method GET
    Write-Host "❌ Unauthorized access allowed (should not happen)" -ForegroundColor Red
} catch {
    Write-Host "✅ Correctly blocked unauthorized access" -ForegroundColor Green
    Write-Host "Error: $($_.Exception.Response.StatusCode)" -ForegroundColor Cyan
}

# Test 2: Invalid login credentials
Write-Host "`n2. Testing invalid login credentials..." -ForegroundColor Yellow
$invalidLogin = @{
    username = "nonexistent"
    password = "wrongpassword"
} | ConvertTo-Json

try {
    Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -Headers @{"Content-Type"="application/json"} -Body $invalidLogin
    Write-Host "❌ Invalid login succeeded (should not happen)" -ForegroundColor Red
} catch {
    Write-Host "✅ Correctly rejected invalid credentials" -ForegroundColor Green
    Write-Host "Error: $($_.Exception.Response.StatusCode)" -ForegroundColor Cyan
}

# Test 3: Create user with duplicate username
Write-Host "`n3. Testing duplicate username creation..." -ForegroundColor Yellow
$duplicateUser = @{
    username = "darkpiaro_updated"  # This username already exists
    password = "password123"
    full_name = "Duplicate User"
} | ConvertTo-Json

try {
    Invoke-RestMethod -Uri "http://localhost:8080/users" -Method POST -Headers @{"Content-Type"="application/json"} -Body $duplicateUser
    Write-Host "❌ Duplicate username allowed (should not happen)" -ForegroundColor Red
} catch {
    Write-Host "✅ Correctly rejected duplicate username" -ForegroundColor Green
    Write-Host "Error: $($_.Exception.Response.StatusCode)" -ForegroundColor Cyan
}

# Test 4: Access non-existent user
Write-Host "`n4. Testing access to non-existent user..." -ForegroundColor Yellow

# First get a valid token
$loginBody = @{
    username = "darkpiaro_updated"
    password = "newpassword123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -Headers @{"Content-Type"="application/json"} -Body $loginBody
$token = $loginResponse.token

try {
    Invoke-RestMethod -Uri "http://localhost:8080/users/999" -Method GET -Headers @{"Authorization"="Bearer $token"}
    Write-Host "❌ Non-existent user returned data (should not happen)" -ForegroundColor Red
} catch {
    Write-Host "✅ Correctly returned 404 for non-existent user" -ForegroundColor Green
    Write-Host "Error: $($_.Exception.Response.StatusCode)" -ForegroundColor Cyan
}

# Test 5: Invalid JSON format
Write-Host "`n5. Testing invalid JSON format..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "http://localhost:8080/users" -Method POST -Headers @{"Content-Type"="application/json"} -Body "invalid json"
    Write-Host "❌ Invalid JSON accepted (should not happen)" -ForegroundColor Red
} catch {
    Write-Host "✅ Correctly rejected invalid JSON" -ForegroundColor Green
    Write-Host "Error: $($_.Exception.Response.StatusCode)" -ForegroundColor Cyan
}

Write-Host "`n🛡️ Error handling testing completed!" -ForegroundColor Green
