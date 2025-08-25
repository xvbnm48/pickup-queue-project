package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"pickup-queue/internal/domain"
	"pickup-queue/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageRepository_Create_HappyPath(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPackageRepository(db)

	pkg := &domain.Package{
		ID:         uuid.New(),
		OrderRef:   "TEST-001",
		DriverCode: "DRV-001",
		Status:     domain.StatusWaiting,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock expectations
	mock.ExpectExec("INSERT INTO packages").
		WithArgs(pkg.ID, pkg.OrderRef, pkg.DriverCode, pkg.Status, pkg.CreatedAt, pkg.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.Create(pkg)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPackageRepository_Create_EdgeCase_DatabaseError(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPackageRepository(db)

	pkg := &domain.Package{
		ID:         uuid.New(),
		OrderRef:   "TEST-001",
		DriverCode: "DRV-001",
		Status:     domain.StatusWaiting,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock expectations - simulate database error
	mock.ExpectExec("INSERT INTO packages").
		WithArgs(pkg.ID, pkg.OrderRef, pkg.DriverCode, pkg.Status, pkg.CreatedAt, pkg.UpdatedAt).
		WillReturnError(sql.ErrConnDone)

	// Execute
	err = repo.Create(pkg)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPackageRepository_GetByID_HappyPath(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPackageRepository(db)

	expectedID := uuid.New()
	expectedTime := time.Now()

	// Mock expectations
	rows := sqlmock.NewRows([]string{
		"id", "order_ref", "driver_code", "status", "created_at", "updated_at",
		"picked_up_at", "handed_over_at", "expired_at",
	}).AddRow(
		expectedID, "TEST-001", "DRV-001", domain.StatusWaiting, expectedTime, expectedTime,
		nil, nil, nil,
	)

	mock.ExpectQuery("SELECT (.+) FROM packages WHERE id = \\$1").
		WithArgs(expectedID).
		WillReturnRows(rows)

	// Execute
	pkg, err := repo.GetByID(expectedID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, expectedID, pkg.ID)
	assert.Equal(t, "TEST-001", pkg.OrderRef)
	assert.Equal(t, "DRV-001", pkg.DriverCode)
	assert.Equal(t, domain.StatusWaiting, pkg.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPackageRepository_GetByID_EdgeCase_NotFound(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPackageRepository(db)

	expectedID := uuid.New()

	// Mock expectations - no rows returned
	mock.ExpectQuery("SELECT (.+) FROM packages WHERE id = \\$1").
		WithArgs(expectedID).
		WillReturnError(sql.ErrNoRows)

	// Execute
	pkg, err := repo.GetByID(expectedID)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, pkg)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPackageRepository_UpdateStatus_HappyPath_StatusPicked(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPackageRepository(db)

	packageID := uuid.New()

	// Mock expectations
	mock.ExpectExec("UPDATE packages SET status = \\$2, updated_at = \\$3, picked_up_at = \\$3 WHERE id = \\$1").
		WithArgs(packageID, domain.StatusPicked, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.UpdateStatus(packageID, domain.StatusPicked)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPackageRepository_UpdateStatus_EdgeCase_PackageNotFound(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPackageRepository(db)

	packageID := uuid.New()

	// Mock expectations - no rows affected (package not found)
	mock.ExpectExec("UPDATE packages SET status = \\$2, updated_at = \\$3, picked_up_at = \\$3 WHERE id = \\$1").
		WithArgs(packageID, domain.StatusPicked, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Execute
	err = repo.UpdateStatus(packageID, domain.StatusPicked)

	// Assert - repository doesn't check affected rows, so no error expected
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPackageRepository_GetExpiredPackages_HappyPath(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPackageRepository(db)

	expiredID := uuid.New()
	expiredTime := time.Now().Add(-25 * time.Hour) // Older than 24 hours

	// Mock expectations
	rows := sqlmock.NewRows([]string{
		"id", "order_ref", "driver_code", "status", "created_at", "updated_at",
		"picked_up_at", "handed_over_at", "expired_at",
	}).AddRow(
		expiredID, "EXPIRED-001", "DRV-001", domain.StatusWaiting, expiredTime, expiredTime,
		nil, nil, nil,
	)

	mock.ExpectQuery("SELECT (.+) FROM packages WHERE status IN \\(\\$1, \\$2\\) AND created_at < \\$3").
		WithArgs(domain.StatusWaiting, domain.StatusPicked, sqlmock.AnyArg()).
		WillReturnRows(rows)

	// Execute
	packages, err := repo.GetExpiredPackages()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, packages, 1)
	assert.Equal(t, expiredID, packages[0].ID)
	assert.Equal(t, domain.StatusWaiting, packages[0].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPackageRepository_GetExpiredPackages_EdgeCase_NoExpiredPackages(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPackageRepository(db)

	// Mock expectations - no rows returned
	rows := sqlmock.NewRows([]string{
		"id", "order_ref", "driver_code", "status", "created_at", "updated_at",
		"picked_up_at", "handed_over_at", "expired_at",
	})

	mock.ExpectQuery("SELECT (.+) FROM packages WHERE status IN \\(\\$1, \\$2\\) AND created_at < \\$3").
		WithArgs(domain.StatusWaiting, domain.StatusPicked, sqlmock.AnyArg()).
		WillReturnRows(rows)

	// Execute
	packages, err := repo.GetExpiredPackages()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, packages, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}
