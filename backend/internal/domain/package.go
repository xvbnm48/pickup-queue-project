package domain

import (
	"time"

	"github.com/google/uuid"
)

// PackageStatus represents the status of a package
type PackageStatus string

const (
	StatusWaiting    PackageStatus = "WAITING"
	StatusPicked     PackageStatus = "PICKED"
	StatusHandedOver PackageStatus = "HANDED_OVER"
	StatusExpired    PackageStatus = "EXPIRED"
)

// Package represents a package in the pickup queue
type Package struct {
	ID           uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderRef     string        `json:"order_reference" gorm:"uniqueIndex;not null"`
	DriverCode   string        `json:"driver_code"`
	Status       PackageStatus `json:"status" gorm:"default:'WAITING'"`
	CreatedAt    time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	PickedUpAt   *time.Time    `json:"picked_up_at,omitempty"`
	HandedOverAt *time.Time    `json:"handed_over_at,omitempty"`
	ExpiredAt    *time.Time    `json:"expired_at,omitempty"`
}

// PackageRepository defines the interface for package data operations
type PackageRepository interface {
	Create(pkg *Package) error
	GetByID(id uuid.UUID) (*Package, error)
	GetByOrderRef(orderRef string) (*Package, error)
	GetAll(limit, offset int, status *PackageStatus) ([]*Package, error)
	Update(pkg *Package) error
	Delete(id uuid.UUID) error
	GetExpiredPackages() ([]*Package, error)
	UpdateStatus(id uuid.UUID, status PackageStatus) error
	GetPackageStats() (*PackageStats, error)
}

// PackageStats represents aggregated package statistics
type PackageStats struct {
	Total      int64 `json:"total"`
	Waiting    int64 `json:"waiting"`
	Picked     int64 `json:"picked"`
	HandedOver int64 `json:"handed_over"`
	Expired    int64 `json:"expired"`
}

// CreatePackageRequest represents the request to create a new package
type CreatePackageRequest struct {
	OrderRef   string `json:"order_reference" binding:"required"`
	DriverCode string `json:"driver_code"`
}

// UpdatePackageStatusRequest represents the request to update package status
type UpdatePackageStatusRequest struct {
	Status PackageStatus `json:"status" binding:"required"`
}
