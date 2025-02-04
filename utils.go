// db/utils.go
package db

import (
	"database/sql"
	"fmt"
	"time"
)

// QueryOptions represents common query parameters
type QueryOptions struct {
	Limit  int
	Offset int
	SortBy string
	Order  string // "ASC" or "DESC"
}

// HealthCheck checks database connectivity and returns status
func HealthCheck(db *sql.DB) error {
	var now time.Time
	err := db.QueryRow("SELECT NOW()").Scan(&now)
	if err != nil {
		return fmt.Errorf("database health check failed: %v", err)
	}
	return nil
}

// GetRecordCount returns total number of records
func GetRecordCount(db *sql.DB) (int64, error) {
	var count int64
	err := db.QueryRow("SELECT COUNT(*) FROM records").Scan(&count)
	return count, err
}

// DeleteRecord deletes a record by ID
func DeleteRecord(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM records WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record with ID %d not found", id)
	}

	return nil
}

// UpdateRecord updates a record
func UpdateRecord(db *sql.DB, record *Record) error {
	query := `
	UPDATE records 
	SET name = $1, 
		description = $2, 
		amount = $3, 
		is_active = $4, 
		updated_at = $5
	WHERE id = $6`

	result, err := db.Exec(query,
		record.Name,
		record.Description,
		record.Amount,
		record.IsActive,
		record.UpdatedAtUnix,
		record.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record with ID %d not found", record.ID)
	}

	return nil
}

// GetRecords retrieves multiple records with pagination and sorting
func GetRecords(db *sql.DB, opts QueryOptions) ([]*Record, error) {
	// Set default values
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.SortBy == "" {
		opts.SortBy = "created_at"
	}
	if opts.Order == "" {
		opts.Order = "DESC"
	}

	query := fmt.Sprintf(`
		SELECT id, name, description, amount, is_active, created_at, updated_at
		FROM records
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, opts.SortBy, opts.Order)

	rows, err := db.Query(query, opts.Limit, opts.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Record
	for rows.Next() {
		record := &Record{}
		err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Description,
			&record.Amount,
			&record.IsActive,
			&record.CreatedAtUnix,
			&record.UpdatedAtUnix,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

// SearchRecords searches records by name or description
func SearchRecords(db *sql.DB, searchTerm string, opts QueryOptions) ([]*Record, error) {
	query := `
		SELECT id, name, description, amount, is_active, created_at, updated_at
		FROM records
		WHERE name ILIKE $1 OR description ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	searchPattern := "%" + searchTerm + "%"

	rows, err := db.Query(query, searchPattern, opts.Limit, opts.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Record
	for rows.Next() {
		record := &Record{}
		err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Description,
			&record.Amount,
			&record.IsActive,
			&record.CreatedAtUnix,
			&record.UpdatedAtUnix,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

// TruncateTable removes all records from the table
func TruncateTable(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE records RESTART IDENTITY")
	return err
}

// CreateIndex creates an index on specified columns
func CreateIndex(db *sql.DB, indexName string, columns []string) error {
	query := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON records (%s)",
		indexName,
		// Join columns with commas
		fmt.Sprintf("%s", columns),
	)
	_, err := db.Exec(query)
	return err
}

// Backup performs a simple backup by copying all records
// Note: This is a basic implementation. For production, use pg_dump
func BackupRecords(db *sql.DB, backupTable string) error {
	queries := []string{
		fmt.Sprintf("DROP TABLE IF EXISTS %s", backupTable),
		fmt.Sprintf("CREATE TABLE %s AS SELECT * FROM records", backupTable),
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
