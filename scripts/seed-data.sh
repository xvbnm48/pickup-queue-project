#!/bin/bash

# Script to create sample packages for testing
API_BASE_URL="http://localhost:8080/api/v1"

echo "🚀 Creating sample packages for testing..."
echo ""

# Sample packages data
packages=(
  '{"order_ref": "ORD-20250824-001", "driver_code": "DRV-JAKARTA-01"}'
  '{"order_ref": "ORD-20250824-002", "driver_code": "DRV-BANDUNG-01"}'
  '{"order_ref": "ORD-20250824-003", "driver_code": "DRV-SURABAYA-01"}'
  '{"order_ref": "ORD-20250824-004", "driver_code": "DRV-MEDAN-01"}'
  '{"order_ref": "ORD-20250824-005", "driver_code": "DRV-JAKARTA-02"}'
  '{"order_ref": "ORD-20250824-006", "driver_code": "DRV-YOGYA-01"}'
  '{"order_ref": "ORD-20250824-007", "driver_code": "DRV-SEMARANG-01"}'
  '{"order_ref": "ORD-20250824-008", "driver_code": "DRV-MALANG-01"}'
  '{"order_ref": "ORD-20250824-009", "driver_code": "DRV-SOLO-01"}'
  '{"order_ref": "ORD-20250824-010", "driver_code": "DRV-DENPASAR-01"}'
)

# Check if API is running
echo "🔍 Checking if API is running..."
if ! curl -s "${API_BASE_URL}/health" > /dev/null; then
  echo "❌ API is not running. Please start the backend server first:"
  echo "   cd backend && go run cmd/api/main.go"
  exit 1
fi

echo "✅ API is running!"
echo ""

# Create packages
counter=1
for package in "${packages[@]}"; do
  echo "📦 Creating package ${counter}/10..."
  
  response=$(curl -s -X POST "${API_BASE_URL}/packages" \
    -H "Content-Type: application/json" \
    -d "$package")
  
  # Extract order_ref from package data
  order_ref=$(echo "$package" | grep -o '"order_ref": "[^"]*"' | cut -d'"' -f4)
  
  if echo "$response" | grep -q '"id"'; then
    echo "   ✅ Created: $order_ref"
  else
    echo "   ❌ Failed: $order_ref"
    echo "   Error: $response"
  fi
  
  ((counter++))
done

echo ""
echo "🎉 Sample data creation completed!"
echo ""
echo "📊 You can now:"
echo "   • View packages: curl ${API_BASE_URL}/packages"
echo "   • Check stats:   curl ${API_BASE_URL}/packages/stats"
echo "   • Open frontend: http://localhost:5173"
