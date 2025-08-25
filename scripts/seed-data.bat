@echo off
:: Script to create sample packages for testing
setlocal enabledelayedexpansion

set "API_BASE_URL=http://localhost:8080/api/v1"

echo 🚀 Creating sample packages for testing...
echo.

:: Check if API is running
echo 🔍 Checking if API is running...
curl -s "%API_BASE_URL%/health" >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ API is not running. Please start the backend server first:
    echo    cd backend ^&^& go run cmd/api/main.go
    exit /b 1
)

echo ✅ API is running!
echo.

:: Create sample packages
set counter=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-001\", \"driver_code\": \"DRV-JAKARTA-01\"}" && echo    ✅ Created: ORD-20250824-001 || echo    ❌ Failed: ORD-20250824-001
set /a counter+=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-002\", \"driver_code\": \"DRV-BANDUNG-01\"}" && echo    ✅ Created: ORD-20250824-002 || echo    ❌ Failed: ORD-20250824-002
set /a counter+=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-003\", \"driver_code\": \"DRV-SURABAYA-01\"}" && echo    ✅ Created: ORD-20250824-003 || echo    ❌ Failed: ORD-20250824-003
set /a counter+=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-004\", \"driver_code\": \"DRV-MEDAN-01\"}" && echo    ✅ Created: ORD-20250824-004 || echo    ❌ Failed: ORD-20250824-004
set /a counter+=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-005\", \"driver_code\": \"DRV-JAKARTA-02\"}" && echo    ✅ Created: ORD-20250824-005 || echo    ❌ Failed: ORD-20250824-005
set /a counter+=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-006\", \"driver_code\": \"DRV-YOGYA-01\"}" && echo    ✅ Created: ORD-20250824-006 || echo    ❌ Failed: ORD-20250824-006
set /a counter+=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-007\", \"driver_code\": \"DRV-SEMARANG-01\"}" && echo    ✅ Created: ORD-20250824-007 || echo    ❌ Failed: ORD-20250824-007
set /a counter+=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-008\", \"driver_code\": \"DRV-MALANG-01\"}" && echo    ✅ Created: ORD-20250824-008 || echo    ❌ Failed: ORD-20250824-008
set /a counter+=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-009\", \"driver_code\": \"DRV-SOLO-01\"}" && echo    ✅ Created: ORD-20250824-009 || echo    ❌ Failed: ORD-20250824-009
set /a counter+=1

echo 📦 Creating package !counter!/10...
curl -s -X POST "%API_BASE_URL%/packages" -H "Content-Type: application/json" -d "{\"order_ref\": \"ORD-20250824-010\", \"driver_code\": \"DRV-DENPASAR-01\"}" && echo    ✅ Created: ORD-20250824-010 || echo    ❌ Failed: ORD-20250824-010

echo.
echo 🎉 Sample data creation completed!
echo.
echo 📊 You can now:
echo    • View packages: curl %API_BASE_URL%/packages
echo    • Check stats:   curl %API_BASE_URL%/packages/stats
echo    • Open frontend: http://localhost:5173
