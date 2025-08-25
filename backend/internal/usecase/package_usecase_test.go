package usecase_test

import (
	"errors"
	"testing"
	"time"

	"pickup-queue/internal/domain"
	"pickup-queue/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
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

func TestPackageUsecase_CreatePackage_HappyPath(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	req := &domain.CreatePackageRequest{
		OrderRef:   "TEST-001",
		DriverCode: "DRV-001",
	}

	// Mock expectations
	mockRepo.On("GetByOrderRef", req.OrderRef).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.Package")).Return(nil)

	// Execute
	pkg, err := uc.CreatePackage(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, req.OrderRef, pkg.OrderRef)
	assert.Equal(t, req.DriverCode, pkg.DriverCode)
	assert.Equal(t, domain.StatusWaiting, pkg.Status)
	assert.NotEqual(t, uuid.Nil, pkg.ID)
	mockRepo.AssertExpectations(t)
}

func TestPackageUsecase_CreatePackage_EdgeCase_DuplicateOrderRef(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	req := &domain.CreatePackageRequest{
		OrderRef:   "DUPLICATE-001",
		DriverCode: "DRV-001",
	}

	existingPkg := &domain.Package{
		ID:       uuid.New(),
		OrderRef: req.OrderRef,
		Status:   domain.StatusWaiting,
	}

	// Mock expectations
	mockRepo.On("GetByOrderRef", req.OrderRef).Return(existingPkg, nil)

	// Execute
	pkg, err := uc.CreatePackage(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pkg)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestPackageUsecase_CreatePackage_EdgeCase_EmptyOrderRef(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	req := &domain.CreatePackageRequest{
		OrderRef:   "", // Empty order reference
		DriverCode: "DRV-001",
	}

	// Execute
	pkg, err := uc.CreatePackage(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pkg)
	assert.Contains(t, err.Error(), "order reference is required")
	mockRepo.AssertNotCalled(t, "GetByOrderRef")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestPackageUsecase_GetPackage_HappyPath(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	packageID := uuid.New()
	expectedPkg := &domain.Package{
		ID:         packageID,
		OrderRef:   "TEST-001",
		DriverCode: "DRV-001",
		Status:     domain.StatusWaiting,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock expectations
	mockRepo.On("GetByID", packageID).Return(expectedPkg, nil)

	// Execute
	pkg, err := uc.GetPackage(packageID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, expectedPkg.ID, pkg.ID)
	assert.Equal(t, expectedPkg.OrderRef, pkg.OrderRef)
	mockRepo.AssertExpectations(t)
}

func TestPackageUsecase_GetPackage_EdgeCase_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	packageID := uuid.New()

	// Mock expectations
	mockRepo.On("GetByID", packageID).Return(nil, nil)

	// Execute
	pkg, err := uc.GetPackage(packageID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pkg)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

func TestPackageUsecase_UpdatePackageStatus_HappyPath_ValidTransition(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	packageID := uuid.New()
	existingPkg := &domain.Package{
		ID:         packageID,
		OrderRef:   "TEST-001",
		DriverCode: "DRV-001",
		Status:     domain.StatusWaiting,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock expectations
	mockRepo.On("GetByID", packageID).Return(existingPkg, nil)
	mockRepo.On("Update", mock.MatchedBy(func(p *domain.Package) bool {
		return p.ID == packageID && p.Status == domain.StatusPicked
	})).Return(nil)

	// Execute
	updatedPkg, err := uc.UpdatePackageStatus(packageID, domain.StatusPicked)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedPkg)
	assert.Equal(t, domain.StatusPicked, updatedPkg.Status)
	mockRepo.AssertExpectations(t)
}

func TestPackageUsecase_UpdatePackageStatus_EdgeCase_InvalidTransition(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	packageID := uuid.New()
	existingPkg := &domain.Package{
		ID:         packageID,
		OrderRef:   "TEST-001",
		DriverCode: "DRV-001",
		Status:     domain.StatusExpired, // Already expired
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock expectations
	mockRepo.On("GetByID", packageID).Return(existingPkg, nil)

	// Execute
	_, err := uc.UpdatePackageStatus(packageID, domain.StatusPicked)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status transition")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestPackageUsecase_MarkExpiredPackages_HappyPath(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	expiredPackages := []*domain.Package{
		{
			ID:        uuid.New(),
			OrderRef:  "EXPIRED-001",
			Status:    domain.StatusWaiting,
			CreatedAt: time.Now().Add(-25 * time.Hour),
		},
		{
			ID:        uuid.New(),
			OrderRef:  "EXPIRED-002",
			Status:    domain.StatusPicked,
			CreatedAt: time.Now().Add(-26 * time.Hour),
		},
	}

	// Mock expectations
	mockRepo.On("GetExpiredPackages").Return(expiredPackages, nil)
	for _, pkg := range expiredPackages {
		mockRepo.On("GetByID", pkg.ID).Return(pkg, nil)
		mockRepo.On("Update", mock.MatchedBy(func(p *domain.Package) bool {
			return p.ID == pkg.ID && p.Status == domain.StatusExpired
		})).Return(nil)
	}

	// Execute
	err := uc.MarkExpiredPackages()

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPackageUsecase_MarkExpiredPackages_EdgeCase_NoExpiredPackages(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	// Mock expectations
	mockRepo.On("GetExpiredPackages").Return([]*domain.Package{}, nil)

	// Execute
	err := uc.MarkExpiredPackages()

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestPackageUsecase_MarkExpiredPackages_EdgeCase_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	uc := usecase.NewPackageUsecase(mockRepo)

	// Mock expectations
	mockRepo.On("GetExpiredPackages").Return([]*domain.Package{}, errors.New("database connection failed"))

	// Execute
	err := uc.MarkExpiredPackages()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection failed")
	mockRepo.AssertExpectations(t)
}
