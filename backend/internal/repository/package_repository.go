package repository

import (
	"pickup-queue/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PackageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) domain.PackageRepository {
	return &PackageRepository{db: db}
}

func (pr *PackageRepository) Create(pkg *domain.Package) error {
	return pr.db.Create(pkg).Error
}

func (pr *PackageRepository) GetByID(id uuid.UUID) (*domain.Package, error) {
	var pkg domain.Package
	err := pr.db.Where("id = ?", id).First(&pkg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pkg, nil
}

func (pr *PackageRepository) GetByOrderRef(orderRef string) (*domain.Package, error) {
	var pkg domain.Package
	err := pr.db.Where("order_ref = ?", orderRef).First(&pkg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pkg, nil
}

func (pr *PackageRepository) GetAll(limit, offset int, status *domain.PackageStatus) ([]*domain.Package, error) {
	var packages []*domain.Package
	query := pr.db.Limit(limit).Offset(offset).Order("created_at DESC")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Find(&packages).Error
	return packages, err
}

func (pr *PackageRepository) Update(pkg *domain.Package) error {
	return pr.db.Save(pkg).Error
}

func (pr *PackageRepository) Delete(id uuid.UUID) error {
	return pr.db.Delete(&domain.Package{}, "id = ?", id).Error
}

func (pr *PackageRepository) GetExpiredPackages() ([]*domain.Package, error) {
	var packages []*domain.Package
	// Packages that have been waiting for more than 24 hours are considered expired
	cutoffTime := time.Now().Add(-24 * time.Hour)

	err := pr.db.Where("status = ? AND created_at < ?", domain.StatusWaiting, cutoffTime).Find(&packages).Error
	return packages, err
}

func (pr *PackageRepository) UpdateStatus(id uuid.UUID, status domain.PackageStatus) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}

	switch status {
	case domain.StatusPicked:
		updates["picked_up_at"] = now
	case domain.StatusHandedOver:
		updates["handed_over_at"] = now
	case domain.StatusExpired:
		updates["expired_at"] = now
	}

	return pr.db.Model(&domain.Package{}).Where("id = ?", id).Updates(updates).Error
}

func (pr *PackageRepository) GetPackageStats() (*domain.PackageStats, error) {
	var stats domain.PackageStats

	// Get total count
	err := pr.db.Model(&domain.Package{}).Count(&stats.Total).Error
	if err != nil {
		return nil, err
	}

	// Get counts by status
	err = pr.db.Model(&domain.Package{}).Where("status = ?", domain.StatusWaiting).Count(&stats.Waiting).Error
	if err != nil {
		return nil, err
	}

	err = pr.db.Model(&domain.Package{}).Where("status = ?", domain.StatusPicked).Count(&stats.Picked).Error
	if err != nil {
		return nil, err
	}

	err = pr.db.Model(&domain.Package{}).Where("status = ?", domain.StatusHandedOver).Count(&stats.HandedOver).Error
	if err != nil {
		return nil, err
	}

	err = pr.db.Model(&domain.Package{}).Where("status = ?", domain.StatusExpired).Count(&stats.Expired).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
