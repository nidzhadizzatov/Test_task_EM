# PowerShell script to run the subscription service

Write-Host "Starting Subscription Service..." -ForegroundColor Green

# Check if database is running (optional)
$env:DATABASE_URL = "postgres://user:password@localhost:5432/subscription_service?sslmode=disable"
$env:PORT = "8080"
$env:LOG_LEVEL = "info"

Write-Host "Environment variables set:" -ForegroundColor Yellow
Write-Host "  DATABASE_URL: $env:DATABASE_URL"
Write-Host "  PORT: $env:PORT"
Write-Host "  LOG_LEVEL: $env:LOG_LEVEL"

# Build and run the application
Write-Host "Building application..." -ForegroundColor Blue
go build -o bin/subscription-service ./cmd/server

Write-Host "Starting server..." -ForegroundColor Blue
./bin/subscription-service