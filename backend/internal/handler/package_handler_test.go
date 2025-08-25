package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pickup-queue/internal/domain"
	"pickup-queue/internal/handler"
	"pickup-queue/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repository untuk testing end-to-end
type MockPackageRepository struct {
	mock.Mock
}

func (m *MockPackageRepository) Create(pkg *domain.Package) error {
	args := m.Called(pkg)
	return args.Error(0)
}

func (m *MockPackageRepository) GetByID(id uuid.UUID) (*domain.Package, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Package), args.Error(1)
}

func (m *MockPackageRepository) GetByOrderRef(orderRef string) (*domain.Package, error) {
	args := m.Called(orderRef)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Package), args.Error(1)
}

func (m *MockPackageRepository) GetAll(limit, offset int, status *domain.PackageStatus) ([]*domain.Package, error) {
	args := m.Called(limit, offset, status)
	return args.Get(0).([]*domain.Package), args.Error(1)
}

func (m *MockPackageRepository) Update(pkg *domain.Package) error {
	args := m.Called(pkg)
	return args.Error(0)
}

func (m *MockPackageRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPackageRepository) UpdateStatus(id uuid.UUID, status domain.PackageStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockPackageRepository) GetExpiredPackages() ([]*domain.Package, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Package), args.Error(1)
}

func (m *MockPackageRepository) GetPackageStats() (*domain.PackageStats, error) {
	args := m.Called()
	return args.Get(0).(*domain.PackageStats), args.Error(1)
}

func setupRouterWithMockRepo(mockRepo *MockPackageRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create real usecase with mock repository
	packageUsecase := usecase.NewPackageUsecase(mockRepo)
	packageHandler := handler.NewPackageHandler(packageUsecase)

	api := router.Group("/api/v1")
	{
		api.POST("/packages", packageHandler.CreatePackage)
		api.GET("/packages", packageHandler.ListPackages)
		api.GET("/packages/:id", packageHandler.GetPackage)
		api.PATCH("/packages/:id/status", packageHandler.UpdatePackageStatus)
		api.DELETE("/packages/:id", packageHandler.DeletePackage)
		api.GET("/packages/stats", packageHandler.GetPackageStats)
	}

	return router
}

func TestPackageHandler_CreatePackage_HappyPath(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	router := setupRouterWithMockRepo(mockRepo)

	requestBody := domain.CreatePackageRequest{
		OrderRef:   "TEST-001",
		DriverCode: "DRV-001",
	}

	// Mock expectations
	mockRepo.On("GetByOrderRef", "TEST-001").Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.Package")).Return(nil)

	// Prepare request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/packages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response["data"])

	packageData := response["data"].(map[string]interface{})
	assert.Equal(t, requestBody.OrderRef, packageData["order_reference"])
	assert.Equal(t, requestBody.DriverCode, packageData["driver_code"])
	assert.Equal(t, "WAITING", packageData["status"])

	mockRepo.AssertExpectations(t)
}

func TestPackageHandler_CreatePackage_EdgeCase_InvalidJSON(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	router := setupRouterWithMockRepo(mockRepo)

	// Prepare invalid JSON request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/packages", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestPackageHandler_GetPackage_HappyPath(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	router := setupRouterWithMockRepo(mockRepo)

	packageID := uuid.New()
	expectedPackage := &domain.Package{
		ID:         packageID,
		OrderRef:   "TEST-001",
		DriverCode: "DRV-001",
		Status:     domain.StatusWaiting,
	}

	// Mock expectations
	mockRepo.On("GetByID", packageID).Return(expectedPackage, nil)

	// Prepare request
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/packages/"+packageID.String(), nil)

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response["data"])

	packageData := response["data"].(map[string]interface{})
	assert.Equal(t, expectedPackage.OrderRef, packageData["order_reference"])
	assert.Equal(t, expectedPackage.DriverCode, packageData["driver_code"])

	mockRepo.AssertExpectations(t)
}

func TestPackageHandler_GetPackageStats_HappyPath(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	router := setupRouterWithMockRepo(mockRepo)

	expectedStats := &domain.PackageStats{
		Total:      100,
		Waiting:    30,
		Picked:     45,
		HandedOver: 20,
		Expired:    5,
	}

	// Mock expectations
	mockRepo.On("GetPackageStats").Return(expectedStats, nil)

	// Prepare request
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/packages/stats", nil)

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response["data"])

	stats := response["data"].(map[string]interface{})
	assert.Equal(t, float64(expectedStats.Total), stats["total"])
	assert.Equal(t, float64(expectedStats.Waiting), stats["waiting"])
	assert.Equal(t, float64(expectedStats.Picked), stats["picked"])

	mockRepo.AssertExpectations(t)
}
