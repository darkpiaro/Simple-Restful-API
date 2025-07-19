# Extended API Testing - Testing UPDATE and DELETE operations

Write-Host "Extended API Testing - CRUD Operations" -ForegroundColor Green

# First, login to get token
Write-Host "`n1. Getting JWT token..." -ForegroundColor Yellow
$loginBody = @{
    username = "darkpiaro"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -Headers @{"Content-Type"="application/json"} -Body $loginBody
    $token = $loginResponse.token
    Write-Host "✅ Token obtained!" -ForegroundColor Green
} catch {
    Write-Host "❌ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    return
}

# Test UPDATE operation
Write-Host "`n2. Testing UPDATE user..." -ForegroundColor Yellow
$updateBody = @{
    username = "darkpiaro_updated"
    full_name = "Dark Piaro Updated"
    password = "newpassword123"
} | ConvertTo-Json

try {
    $updateResponse = Invoke-RestMethod -Uri "http://localhost:8080/users/1" -Method PUT -Headers @{"Content-Type"="application/json"; "Authorization"="Bearer $token"} -Body $updateBody
    Write-Host "✅ User updated successfully!" -ForegroundColor Green
    $updateResponse | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Update failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test login with new password
Write-Host "`n3. Testing login with updated password..." -ForegroundColor Yellow
$newLoginBody = @{
    username = "darkpiaro_updated"
    password = "newpassword123"
} | ConvertTo-Json

try {
    $newLoginResponse = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -Headers @{"Content-Type"="application/json"} -Body $newLoginBody
    Write-Host "✅ Login with new password successful!" -ForegroundColor Green
    $newToken = $newLoginResponse.token
} catch {
    Write-Host "❌ Login with new password failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Create another user for testing
Write-Host "`n4. Creating second user for deletion test..." -ForegroundColor Yellow
$body2 = @{
    username = "testuser2"
    password = "password123"
    full_name = "Test User 2"
} | ConvertTo-Json

try {
    $response2 = Invoke-RestMethod -Uri "http://localhost:8080/users" -Method POST -Headers @{"Content-Type"="application/json"} -Body $body2
    Write-Host "✅ Second user created!" -ForegroundColor Green
    $userId2 = $response2.user.id
} catch {
    Write-Host "❌ Failed to create second user: $($_.Exception.Message)" -ForegroundColor Red
    return
}

# Test DELETE operation
Write-Host "`n5. Testing DELETE user..." -ForegroundColor Yellow
try {
    $deleteResponse = Invoke-RestMethod -Uri "http://localhost:8080/users/$userId2" -Method DELETE -Headers @{"Authorization"="Bearer $newToken"}
    Write-Host "✅ User deleted successfully!" -ForegroundColor Green
    $deleteResponse | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Delete failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Verify user was deleted
Write-Host "`n6. Verifying user deletion..." -ForegroundColor Yellow
try {
    $verifyResponse = Invoke-RestMethod -Uri "http://localhost:8080/users/$userId2" -Method GET -Headers @{"Authorization"="Bearer $newToken"}
    Write-Host "❌ User still exists (deletion failed)" -ForegroundColor Red
} catch {
    Write-Host "✅ User successfully deleted (404 error expected)" -ForegroundColor Green
}

# Final check - get all users
Write-Host "`n7. Final check - all remaining users..." -ForegroundColor Yellow
try {
    $finalUsers = Invoke-RestMethod -Uri "http://localhost:8080/users" -Method GET -Headers @{"Authorization"="Bearer $newToken"}
    Write-Host "✅ Retrieved final user list!" -ForegroundColor Green
    $finalUsers | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Failed to get final user list: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🎉 Extended CRUD testing completed!" -ForegroundColor Green
