# PowerShell script to create sample packages for testing

$API_BASE_URL = "http://localhost:8080/api/v1"

Write-Host "🚀 Creating sample packages for testing..." -ForegroundColor Green
Write-Host ""

# Sample packages data
$packages = @(
    @{ order_ref = "ORD-20250824-001"; driver_code = "DRV-JAKARTA-01" },
    @{ order_ref = "ORD-20250824-002"; driver_code = "DRV-BANDUNG-01" },
    @{ order_ref = "ORD-20250824-003"; driver_code = "DRV-SURABAYA-01" },
    @{ order_ref = "ORD-20250824-004"; driver_code = "DRV-MEDAN-01" },
    @{ order_ref = "ORD-20250824-005"; driver_code = "DRV-JAKARTA-02" },
    @{ order_ref = "ORD-20250824-006"; driver_code = "DRV-YOGYA-01" },
    @{ order_ref = "ORD-20250824-007"; driver_code = "DRV-SEMARANG-01" },
    @{ order_ref = "ORD-20250824-008"; driver_code = "DRV-MALANG-01" },
    @{ order_ref = "ORD-20250824-009"; driver_code = "DRV-SOLO-01" },
    @{ order_ref = "ORD-20250824-010"; driver_code = "DRV-DENPASAR-01" }
)

# Check if API is running
Write-Host "🔍 Checking if API is running..." -ForegroundColor Yellow
try {
    $healthCheck = Invoke-RestMethod -Uri "$API_BASE_URL/health" -Method Get -TimeoutSec 5
    Write-Host "✅ API is running!" -ForegroundColor Green
}
catch {
    Write-Host "❌ API is not running. Please start the backend server first:" -ForegroundColor Red
    Write-Host "   cd backend && go run cmd/api/main.go" -ForegroundColor White
    exit 1
}

Write-Host ""

# Create packages
$counter = 1
foreach ($package in $packages) {
    Write-Host "📦 Creating package $counter/10..." -ForegroundColor Cyan
    
    $body = @{
        order_ref = $package.order_ref
        driver_code = $package.driver_code
    } | ConvertTo-Json
    
    try {
        $response = Invoke-RestMethod -Uri "$API_BASE_URL/packages" -Method Post -Body $body -ContentType "application/json"
        Write-Host "   ✅ Created: $($package.order_ref)" -ForegroundColor Green
    }
    catch {
        Write-Host "   ❌ Failed: $($package.order_ref)" -ForegroundColor Red
        Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    $counter++
}

Write-Host ""
Write-Host "🎉 Sample data creation completed!" -ForegroundColor Green
Write-Host ""
Write-Host "📊 You can now:" -ForegroundColor Cyan
Write-Host "   • View packages: curl $API_BASE_URL/packages" -ForegroundColor White
Write-Host "   • Check stats:   curl $API_BASE_URL/packages/stats" -ForegroundColor White
Write-Host "   • Open frontend: http://localhost:5173" -ForegroundColor White
