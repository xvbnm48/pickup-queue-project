package repository

import (
	"database/sql"
	"fmt"
	"pickup-queue/internal/domain"
	"pickup-queue/pkg/database"
	"time"

	"github.com/google/uuid"
)

type PackageRepository struct {
	db *sql.DB
}

func NewPackageRepository(db *sql.DB) domain.PackageRepository {
	return &PackageRepository{db: db}
}

func (pr *PackageRepository) Create(pkg *domain.Package) error {
	query := `
		INSERT INTO packages (id, order_ref, driver_code, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	args := []interface{}{
		pkg.ID,
		pkg.OrderRef,
		pkg.DriverCode,
		pkg.Status,
		pkg.CreatedAt,
		pkg.UpdatedAt,
	}

	startTime := time.Now()
	_, err := pr.db.Exec(query, args...)

	if err != nil {
		database.LogQueryError(query, args, err, startTime)
	} else {
		database.LogQuery(query, args, startTime)
	}

	return err
}

func (pr *PackageRepository) GetByID(id uuid.UUID) (*domain.Package, error) {
	query := `
		SELECT id, order_ref, driver_code, status, created_at, updated_at, 
		       picked_up_at, handed_over_at, expired_at
		FROM packages 
		WHERE id = $1`

	args := []interface{}{id}
	startTime := time.Now()

	var pkg domain.Package
	var pickedUpAt, handedOverAt, expiredAt sql.NullTime

	err := pr.db.QueryRow(query, id).Scan(
		&pkg.ID,
		&pkg.OrderRef,
		&pkg.DriverCode,
		&pkg.Status,
		&pkg.CreatedAt,
		&pkg.UpdatedAt,
		&pickedUpAt,
		&handedOverAt,
		&expiredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			database.LogQuery(query, args, startTime)
			return nil, nil
		}
		database.LogQueryError(query, args, err, startTime)
		return nil, err
	}

	database.LogQuery(query, args, startTime)

	// Handle nullable time fields
	if pickedUpAt.Valid {
		pkg.PickedUpAt = &pickedUpAt.Time
	}
	if handedOverAt.Valid {
		pkg.HandedOverAt = &handedOverAt.Time
	}
	if expiredAt.Valid {
		pkg.ExpiredAt = &expiredAt.Time
	}

	return &pkg, nil
}

func (pr *PackageRepository) GetByOrderRef(orderRef string) (*domain.Package, error) {
	query := `
		SELECT id, order_ref, driver_code, status, created_at, updated_at, 
		       picked_up_at, handed_over_at, expired_at
		FROM packages 
		WHERE order_ref = $1`

	args := []interface{}{orderRef}
	startTime := time.Now()

	var pkg domain.Package
	var pickedUpAt, handedOverAt, expiredAt sql.NullTime

	err := pr.db.QueryRow(query, orderRef).Scan(
		&pkg.ID,
		&pkg.OrderRef,
		&pkg.DriverCode,
		&pkg.Status,
		&pkg.CreatedAt,
		&pkg.UpdatedAt,
		&pickedUpAt,
		&handedOverAt,
		&expiredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			database.LogQuery(query, args, startTime)
			return nil, nil
		}
		database.LogQueryError(query, args, err, startTime)
		return nil, err
	}

	database.LogQuery(query, args, startTime)

	// Handle nullable time fields
	if pickedUpAt.Valid {
		pkg.PickedUpAt = &pickedUpAt.Time
	}
	if handedOverAt.Valid {
		pkg.HandedOverAt = &handedOverAt.Time
	}
	if expiredAt.Valid {
		pkg.ExpiredAt = &expiredAt.Time
	}

	return &pkg, nil
}

func (pr *PackageRepository) GetAll(limit, offset int, status *domain.PackageStatus) ([]*domain.Package, error) {
	baseQuery := `
		SELECT id, order_ref, driver_code, status, created_at, updated_at, 
		       picked_up_at, handed_over_at, expired_at
		FROM packages`

	var args []interface{}
	var whereClause string
	argIndex := 1

	if status != nil {
		whereClause = " WHERE status = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *status)
		argIndex++
	}

	query := baseQuery + whereClause +
		" ORDER BY created_at DESC" +
		" LIMIT $" + fmt.Sprintf("%d", argIndex) +
		" OFFSET $" + fmt.Sprintf("%d", argIndex+1)

	args = append(args, limit, offset)

	startTime := time.Now()
	rows, err := pr.db.Query(query, args...)
	if err != nil {
		database.LogQueryError(query, args, err, startTime)
		return nil, err
	}
	defer rows.Close()

	database.LogQuery(query, args, startTime)

	var packages []*domain.Package
	for rows.Next() {
		var pkg domain.Package
		var pickedUpAt, handedOverAt, expiredAt sql.NullTime

		err := rows.Scan(
			&pkg.ID,
			&pkg.OrderRef,
			&pkg.DriverCode,
			&pkg.Status,
			&pkg.CreatedAt,
			&pkg.UpdatedAt,
			&pickedUpAt,
			&handedOverAt,
			&expiredAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable time fields
		if pickedUpAt.Valid {
			pkg.PickedUpAt = &pickedUpAt.Time
		}
		if handedOverAt.Valid {
			pkg.HandedOverAt = &handedOverAt.Time
		}
		if expiredAt.Valid {
			pkg.ExpiredAt = &expiredAt.Time
		}

		packages = append(packages, &pkg)
	}

	return packages, rows.Err()
}

func (pr *PackageRepository) Update(pkg *domain.Package) error {
	query := `
		UPDATE packages 
		SET order_ref = $2, driver_code = $3, status = $4, updated_at = $5,
		    picked_up_at = $6, handed_over_at = $7, expired_at = $8
		WHERE id = $1`

	args := []interface{}{
		pkg.ID,
		pkg.OrderRef,
		pkg.DriverCode,
		pkg.Status,
		pkg.UpdatedAt,
		pkg.PickedUpAt,
		pkg.HandedOverAt,
		pkg.ExpiredAt,
	}

	startTime := time.Now()
	_, err := pr.db.Exec(query, args...)

	if err != nil {
		database.LogQueryError(query, args, err, startTime)
	} else {
		database.LogQuery(query, args, startTime)
	}

	return err
}

func (pr *PackageRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM packages WHERE id = $1`
	args := []interface{}{id}

	startTime := time.Now()
	_, err := pr.db.Exec(query, id)

	if err != nil {
		database.LogQueryError(query, args, err, startTime)
	} else {
		database.LogQuery(query, args, startTime)
	}

	return err
}

func (pr *PackageRepository) GetExpiredPackages() ([]*domain.Package, error) {
	// Packages that have been waiting for more than 24 hours are considered expired
	cutoffTime := time.Now().Add(-24 * time.Hour)

	query := `
		SELECT id, order_ref, driver_code, status, created_at, updated_at, 
		       picked_up_at, handed_over_at, expired_at
		FROM packages 
		WHERE status IN ($1, $2) AND created_at < $3`

	args := []interface{}{domain.StatusWaiting, domain.StatusPicked, cutoffTime}
	startTime := time.Now()

	rows, err := pr.db.Query(query, domain.StatusWaiting, domain.StatusPicked, cutoffTime)
	if err != nil {
		database.LogQueryError(query, args, err, startTime)
		return nil, err
	}
	defer rows.Close()

	database.LogQuery(query, args, startTime)

	var packages []*domain.Package
	for rows.Next() {
		var pkg domain.Package
		var pickedUpAt, handedOverAt, expiredAt sql.NullTime

		err := rows.Scan(
			&pkg.ID,
			&pkg.OrderRef,
			&pkg.DriverCode,
			&pkg.Status,
			&pkg.CreatedAt,
			&pkg.UpdatedAt,
			&pickedUpAt,
			&handedOverAt,
			&expiredAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable time fields
		if pickedUpAt.Valid {
			pkg.PickedUpAt = &pickedUpAt.Time
		}
		if handedOverAt.Valid {
			pkg.HandedOverAt = &handedOverAt.Time
		}
		if expiredAt.Valid {
			pkg.ExpiredAt = &expiredAt.Time
		}

		packages = append(packages, &pkg)
	}

	return packages, rows.Err()
}

func (pr *PackageRepository) UpdateStatus(id uuid.UUID, status domain.PackageStatus) error {
	now := time.Now()

	var query string
	var args []interface{}

	switch status {
	case domain.StatusPicked:
		query = `UPDATE packages SET status = $2, updated_at = $3, picked_up_at = $3 WHERE id = $1`
		args = []interface{}{id, status, now}
	case domain.StatusHandedOver:
		query = `UPDATE packages SET status = $2, updated_at = $3, handed_over_at = $3 WHERE id = $1`
		args = []interface{}{id, status, now}
	case domain.StatusExpired:
		query = `UPDATE packages SET status = $2, updated_at = $3, expired_at = $3 WHERE id = $1`
		args = []interface{}{id, status, now}
	default:
		query = `UPDATE packages SET status = $2, updated_at = $3 WHERE id = $1`
		args = []interface{}{id, status, now}
	}

	startTime := time.Now()
	_, err := pr.db.Exec(query, args...)

	if err != nil {
		database.LogQueryError(query, args, err, startTime)
	} else {
		database.LogQuery(query, args, startTime)
	}

	return err
}

func (pr *PackageRepository) GetPackageStats() (*domain.PackageStats, error) {
	var stats domain.PackageStats

	// Get total count
	query1 := "SELECT COUNT(*) FROM packages"
	args1 := []interface{}{}
	startTime1 := time.Now()
	err := pr.db.QueryRow(query1).Scan(&stats.Total)
	if err != nil {
		database.LogQueryError(query1, args1, err, startTime1)
		return nil, err
	}
	database.LogQuery(query1, args1, startTime1)

	// Get counts by status
	query2 := "SELECT COUNT(*) FROM packages WHERE status = $1"
	args2 := []interface{}{domain.StatusWaiting}
	startTime2 := time.Now()
	err = pr.db.QueryRow(query2, domain.StatusWaiting).Scan(&stats.Waiting)
	if err != nil {
		database.LogQueryError(query2, args2, err, startTime2)
		return nil, err
	}
	database.LogQuery(query2, args2, startTime2)

	args3 := []interface{}{domain.StatusPicked}
	startTime3 := time.Now()
	err = pr.db.QueryRow(query2, domain.StatusPicked).Scan(&stats.Picked)
	if err != nil {
		database.LogQueryError(query2, args3, err, startTime3)
		return nil, err
	}
	database.LogQuery(query2, args3, startTime3)

	args4 := []interface{}{domain.StatusHandedOver}
	startTime4 := time.Now()
	err = pr.db.QueryRow(query2, domain.StatusHandedOver).Scan(&stats.HandedOver)
	if err != nil {
		database.LogQueryError(query2, args4, err, startTime4)
		return nil, err
	}
	database.LogQuery(query2, args4, startTime4)

	args5 := []interface{}{domain.StatusExpired}
	startTime5 := time.Now()
	err = pr.db.QueryRow(query2, domain.StatusExpired).Scan(&stats.Expired)
	if err != nil {
		database.LogQueryError(query2, args5, err, startTime5)
		return nil, err
	}
	database.LogQuery(query2, args5, startTime5)

	return &stats, nil
}
