package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite" // pure-Go driver, registers as "sqlite" — no cgo/C compiler required
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Open opens (creating if needed) the SQLite database at path, enables WAL mode
// and foreign keys, and runs all pending migrations.
// Uses modernc.org/sqlite (pure Go, no cgo) so no C compiler is needed to build
// on Windows — avoids 32-bit/64-bit MinGW toolchain mismatches entirely.
func Open(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("ouverture base de données: %w", err)
	}
	db.SetMaxOpenConns(1) // SQLite + single desktop user: avoid locked-DB contention

	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrations: %w", err)
	}
	return db, nil
}

func runMigrations(db *sql.DB) error {
	srcDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}
	dbDriver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", srcDriver, "sqlite", dbDriver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Backup copies the live database to destPath using SQLite's online backup
// (VACUUM INTO), which is safe to run while the app is connected.
func Backup(db *sql.DB, destPath string) error {
	_, err := db.Exec("VACUUM INTO ?", destPath)
	if err != nil {
		return fmt.Errorf("échec de la sauvegarde: %w", err)
	}
	return nil
}
