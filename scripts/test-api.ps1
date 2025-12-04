# Test Script for Subscription Service API

Write-Host "Testing Subscription Service API..." -ForegroundColor Green

$BASE_URL = "http://localhost:8080"
$API_URL = "$BASE_URL/api/v1"

# Test data
$testUserID = "60601fee-2bf1-4721-ae6f-7636e79a0cba"
$subscriptionData = @{
    service_name = "Yandex Plus"
    price = 400
    user_id = $testUserID
    start_date = "07-2025"
} | ConvertTo-Json

Write-Host "Testing API endpoints..." -ForegroundColor Blue

# Test 1: Health check
Write-Host "1. Testing health check..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/health" -Method GET
    Write-Host "✓ Health check passed: $($response | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "✗ Health check failed: $_" -ForegroundColor Red
}

# Test 2: Get all subscriptions (should be empty initially)
Write-Host "2. Testing GET all subscriptions..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$API_URL/subscriptions" -Method GET
    Write-Host "✓ GET all subscriptions: Found $($response.Count) subscriptions" -ForegroundColor Green
} catch {
    Write-Host "✗ GET all subscriptions failed: $_" -ForegroundColor Red
}

# Test 3: Create subscription
Write-Host "3. Testing POST create subscription..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$API_URL/subscriptions" -Method POST -Body $subscriptionData -ContentType "application/json"
    $createdId = $response.id
    Write-Host "✓ Created subscription with ID: $createdId" -ForegroundColor Green
} catch {
    Write-Host "✗ Create subscription failed: $_" -ForegroundColor Red
}

# Test 4: Get subscription by ID
if ($createdId) {
    Write-Host "4. Testing GET subscription by ID..." -ForegroundColor Yellow
    try {
        $response = Invoke-RestMethod -Uri "$API_URL/subscriptions/$createdId" -Method GET
        Write-Host "✓ GET subscription by ID successful" -ForegroundColor Green
    } catch {
        Write-Host "✗ GET subscription by ID failed: $_" -ForegroundColor Red
    }
}

# Test 5: Calculate cost
Write-Host "5. Testing cost calculation..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$API_URL/subscriptions/cost?user_id=$testUserID" -Method GET
    Write-Host "✓ Cost calculation: Total cost = $($response.total_cost) rubles" -ForegroundColor Green
} catch {
    Write-Host "✗ Cost calculation failed: $_" -ForegroundColor Red
}

Write-Host "API testing completed!" -ForegroundColor Green