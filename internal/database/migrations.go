package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/001_create_requests_table.sql
var migrationFS embed.FS

type Migration struct {
	Version   string
	Filename  string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
	IsApplied bool
}

type MigrationRunner struct {
	db         *sql.DB
	migrations []Migration
}

func NewMigrationRunner(db *sql.DB, _ string) *MigrationRunner {
	return &MigrationRunner{
		db: db,
	}
}

func (mr *MigrationRunner) LoadMigrations() error {
	if err := mr.loadEmbeddedMigrations(); err != nil {
		return fmt.Errorf("failed to load embedded migrations: %w", err)
	}

	log.Printf("Loaded %d embedded migrations", len(mr.migrations))

	return nil
}

func (mr *MigrationRunner) loadEmbeddedMigrations() error {
	migrationFile := "migrations/001_create_requests_table.sql"

	content, err := migrationFS.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read embedded migration: %w", err)
	}

	migration, err := mr.parseMigration("001_create_requests_table.sql", string(content))
	if err != nil {
		return fmt.Errorf("failed to parse embedded migration: %w", err)
	}

	mr.migrations = append(mr.migrations, migration)

	log.Printf("Successfully loaded embedded migration: %s", migrationFile)
	return nil
}

func (mr *MigrationRunner) parseMigration(filename, content string) (Migration, error) {
	version := strings.Split(filename, "_")[0]

	parts := strings.Split(content, "-- +migrate Down")
	if len(parts) != 2 {
		return Migration{}, fmt.Errorf("migration %s missing -- +migrate Down marker", filename)
	}

	upSQL := strings.TrimSpace(parts[0])
	downSQL := strings.TrimSpace(parts[1])

	upSQL = strings.Replace(upSQL, "-- +migrate Up", "", 1)
	upSQL = strings.TrimSpace(upSQL)

	return Migration{
		Version:  version,
		Filename: filename,
		UpSQL:    upSQL,
		DownSQL:  downSQL,
	}, nil
}

func (mr *MigrationRunner) EnsureMigrationTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		filename TEXT NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := mr.db.Exec(query)
	return err
}

func (mr *MigrationRunner) GetAppliedMigrations() (map[string]bool, error) {
	rows, err := mr.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

func (mr *MigrationRunner) Up() error {
	if err := mr.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	applied, err := mr.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, migration := range mr.migrations {
		if applied[migration.Version] {
			migration.IsApplied = true
			log.Printf("Migration %s already applied", migration.Version)
			continue
		}

		if err := mr.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		log.Printf("Applied migration %s (%s)", migration.Version, migration.Filename)
	}

	return nil
}

func (mr *MigrationRunner) applyMigration(migration Migration) error {
	tx, err := mr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	insertQuery := `
	INSERT INTO schema_migrations (version, filename) 
	VALUES (?, ?)`

	if _, err := tx.Exec(insertQuery, migration.Version, migration.Filename); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

func (mr *MigrationRunner) GetMigrationStatus() ([]Migration, error) {
	if err := mr.EnsureMigrationTable(); err != nil {
		return nil, err
	}

	applied, err := mr.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	for i := range mr.migrations {
		mr.migrations[i].IsApplied = applied[mr.migrations[i].Version]
	}

	return mr.migrations, nil
}
