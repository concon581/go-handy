// db/postgres.go
package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// Record represents a record with Unix timestamps
type Record struct {
	ID          int64    `db:"id"`
	Name        string   `db:"name"`
	Description *string  `db:"description"`
	Amount      *float64 `db:"amount"`
	IsActive    *bool    `db:"is_active"`
	// Unix timestamps as int64
	CreatedAtUnix int64  `db:"created_at"`
	UpdatedAtUnix *int64 `db:"updated_at"` // Nullable
}

// ToTime converts Unix timestamp to time.Time
func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// FromTime converts time.Time to Unix timestamp
func TimeToUnix(t time.Time) int64 {
	return t.Unix()
}

// CreateTable creates the records table with BIGINT for timestamps
func CreateTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS records (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		amount NUMERIC(15,2),
		is_active BOOLEAN,
		created_at BIGINT NOT NULL,    -- Unix timestamp
		updated_at BIGINT              -- Nullable Unix timestamp
	);
	`
	_, err := db.Exec(query)
	return err
}

// InsertRecord inserts a single record, working with Unix timestamps
func InsertRecord(db *sql.DB, record *Record) error {
	query := `
	INSERT INTO records (
		name, description, amount, is_active, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6
	) RETURNING id`

	err := db.QueryRow(
		query,
		record.Name,
		record.Description,
		record.Amount,
		record.IsActive,
		record.CreatedAtUnix, // Store as Unix timestamp
		record.UpdatedAtUnix, // Store as Unix timestamp
	).Scan(&record.ID)

	return err
}

// GetRecord retrieves a record by ID
func GetRecord(db *sql.DB, id int64) (*Record, error) {
	record := &Record{}
	query := `
	SELECT id, name, description, amount, is_active, created_at, updated_at
	FROM records WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
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

	return record, nil
}
