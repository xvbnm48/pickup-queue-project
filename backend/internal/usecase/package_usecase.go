package usecase

import (
	"errors"
	"pickup-queue/internal/domain"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPackageNotFound         = errors.New("package not found")
	ErrDuplicateOrderRef       = errors.New("order reference already exists")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

type PackageUsecase struct {
	packageRepo domain.PackageRepository
}

func NewPackageUsecase(packageRepo domain.PackageRepository) *PackageUsecase {
	return &PackageUsecase{
		packageRepo: packageRepo,
	}
}

func (pu *PackageUsecase) CreatePackage(req *domain.CreatePackageRequest) (*domain.Package, error) {
	// Check if order reference already exists
	existing, _ := pu.packageRepo.GetByOrderRef(req.OrderRef)
	if existing != nil {
		return nil, ErrDuplicateOrderRef
	}

	pkg := &domain.Package{
		ID:         uuid.New(),
		OrderRef:   req.OrderRef,
		DriverCode: req.DriverCode,
		Status:     domain.StatusWaiting,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := pu.packageRepo.Create(pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func (pu *PackageUsecase) GetPackage(id uuid.UUID) (*domain.Package, error) {
	pkg, err := pu.packageRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, ErrPackageNotFound
	}
	return pkg, nil
}

func (pu *PackageUsecase) GetPackageByOrderRef(orderRef string) (*domain.Package, error) {
	pkg, err := pu.packageRepo.GetByOrderRef(orderRef)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, ErrPackageNotFound
	}
	return pkg, nil
}

func (pu *PackageUsecase) ListPackages(limit, offset int, status *domain.PackageStatus) ([]*domain.Package, error) {
	return pu.packageRepo.GetAll(limit, offset, status)
}

func (pu *PackageUsecase) UpdatePackageStatus(id uuid.UUID, newStatus domain.PackageStatus) (*domain.Package, error) {
	pkg, err := pu.packageRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, ErrPackageNotFound
	}

	// Validate status transition
	if !pu.isValidStatusTransition(pkg.Status, newStatus) {
		return nil, ErrInvalidStatusTransition
	}

	now := time.Now()
	pkg.Status = newStatus
	pkg.UpdatedAt = now

	switch newStatus {
	case domain.StatusPicked:
		pkg.PickedUpAt = &now
	case domain.StatusHandedOver:
		pkg.HandedOverAt = &now
	case domain.StatusExpired:
		pkg.ExpiredAt = &now
	}

	err = pu.packageRepo.Update(pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func (pu *PackageUsecase) DeletePackage(id uuid.UUID) error {
	pkg, err := pu.packageRepo.GetByID(id)
	if err != nil {
		return err
	}
	if pkg == nil {
		return ErrPackageNotFound
	}

	return pu.packageRepo.Delete(id)
}

func (pu *PackageUsecase) GetPackageStats() (*domain.PackageStats, error) {
	return pu.packageRepo.GetPackageStats()
}

func (pu *PackageUsecase) MarkExpiredPackages() error {
	expiredPackages, err := pu.packageRepo.GetExpiredPackages()
	if err != nil {
		return err
	}

	for _, pkg := range expiredPackages {
		_, err := pu.UpdatePackageStatus(pkg.ID, domain.StatusExpired)
		if err != nil {
			// Log error but continue processing other packages
			continue
		}
	}

	return nil
}

func (pu *PackageUsecase) isValidStatusTransition(currentStatus, newStatus domain.PackageStatus) bool {
	switch currentStatus {
	case domain.StatusWaiting:
		return newStatus == domain.StatusPicked || newStatus == domain.StatusExpired
	case domain.StatusPicked:
		return newStatus == domain.StatusHandedOver || newStatus == domain.StatusExpired
	case domain.StatusHandedOver, domain.StatusExpired:
		return false // Terminal states
	default:
		return false
	}
}
