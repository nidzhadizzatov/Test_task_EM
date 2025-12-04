# Database Migration Script for Windows

Write-Host "Running Database Migrations..." -ForegroundColor Green

# Database connection parameters from environment variables
$DB_URL = $env:DATABASE_URL
if (-not $DB_URL) {
    $DB_URL = "postgres://user:password@localhost:5432/subscription_service?sslmode=disable"
}

Write-Host "Database URL: $DB_URL" -ForegroundColor Yellow

# Check if psql is available
$psql = Get-Command psql -ErrorAction SilentlyContinue
if (-not $psql) {
    Write-Host "Error: psql command not found. Please install PostgreSQL client tools." -ForegroundColor Red
    Write-Host "You can download PostgreSQL from: https://www.postgresql.org/download/windows/" -ForegroundColor Yellow
    exit 1
}

# Run migrations
Write-Host "Applying migrations..." -ForegroundColor Blue
try {
    & psql $DB_URL -f "db/migrations/0001_init.up.sql"
    Write-Host "Migrations completed successfully!" -ForegroundColor Green
}
catch {
    Write-Host "Error running migrations: $_" -ForegroundColor Red
    Write-Host "Please ensure PostgreSQL is running and connection details are correct." -ForegroundColor Yellow
    exit 1
}