package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CreatePackageRequest struct {
	OrderRef   string `json:"order_reference"`
	DriverCode string `json:"driver_code"`
}

func main() {
	baseURL := "http://localhost:8080"

	// Create test packages
	packages := []CreatePackageRequest{
		{OrderRef: "ABC-001", DriverCode: "DRV-001"},
		{OrderRef: "ABC-002", DriverCode: "DRV-002"},
		{OrderRef: "ABC-003", DriverCode: "DRV-003"},
		{OrderRef: "DEF-001", DriverCode: "DRV-001"},
		{OrderRef: "DEF-002", DriverCode: "DRV-004"},
		{OrderRef: "GHI-001", DriverCode: ""},
		{OrderRef: "GHI-002", DriverCode: "DRV-002"},
		{OrderRef: "JKL-001", DriverCode: "DRV-005"},
	}

	fmt.Println("Creating test packages...")

	for _, pkg := range packages {
		// Create package
		jsonData, _ := json.Marshal(pkg)
		resp, err := http.Post(baseURL+"/api/v1/packages", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error creating package %s: %v", pkg.OrderRef, err)
			continue
		}

		if resp.StatusCode == 201 || resp.StatusCode == 409 { // 409 means already exists
			fmt.Printf("✓ Package %s created successfully\n", pkg.OrderRef)
		} else {
			fmt.Printf("✗ Failed to create package %s (status: %d)\n", pkg.OrderRef, resp.StatusCode)
		}
		resp.Body.Close()

		// Add some delay
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nAll test packages have been created!")
	fmt.Println("You can now test the frontend at http://localhost:3000")
}
