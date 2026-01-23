package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var migrationFS embed.FS

type Migration struct {
	Version    string
	Filename   string
	UpSQL      string
	DownSQL    string
	AppliedAt  *time.Time
	IsApplied  bool
	IsEmbedded bool
}

type MigrationRunner struct {
	db           *sql.DB
	migrations   []Migration
	useEmbedded  bool
	externalPath string
}

func NewMigrationRunner(db *sql.DB, externalPath string) *MigrationRunner {
	return &MigrationRunner{
		db:           db,
		externalPath: externalPath,
		useEmbedded:  externalPath == "",
	}
}

func (mr *MigrationRunner) LoadMigrations() error {
	if err := mr.loadEmbeddedMigrations(); err != nil {
		return fmt.Errorf("failed to load embedded migrations: %w", err)
	}

	if mr.externalPath != "" {
		if _, err := os.Stat(mr.externalPath); err == nil {
			if err := mr.loadExternalMigrations(); err != nil {
				log.Printf("Warning: failed to load external migrations: %v", err)
			}
		}
	}

	sort.Slice(mr.migrations, func(i, j int) bool {
		return mr.migrations[i].Version < mr.migrations[j].Version
	})

	log.Printf("Loaded %d migrations (embedded: %t, external: %s)",
		len(mr.migrations), mr.useEmbedded, mr.externalPath)

	return nil
}

func (mr *MigrationRunner) loadEmbeddedMigrations() error {
	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read embedded migrations: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := migrationFS.ReadFile(filepath.Join("migrations", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read embedded migration %s: %w", entry.Name(), err)
		}

		migration, err := mr.parseMigration(entry.Name(), string(content))
		if err != nil {
			return fmt.Errorf("failed to parse embedded migration %s: %w", entry.Name(), err)
		}

		migration.IsEmbedded = true
		mr.migrations = append(mr.migrations, migration)
	}

	return nil
}

func (mr *MigrationRunner) loadExternalMigrations() error {
	entries, err := os.ReadDir(mr.externalPath)
	if err != nil {
		return fmt.Errorf("failed to read external migrations: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(mr.externalPath, entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read external migration %s: %w", entry.Name(), err)
		}

		migration, err := mr.parseMigration(entry.Name(), string(content))
		if err != nil {
			return fmt.Errorf("failed to parse external migration %s: %w", entry.Name(), err)
		}

		migration.IsEmbedded = false
		mr.migrations = append(mr.migrations, migration)
	}

	return nil
}

func (mr *MigrationRunner) parseMigration(filename, content string) (Migration, error) {
	version := strings.Split(filename, "_")[0]

	parts := strings.Split(content, "-- +goose Down")
	if len(parts) != 2 {
		return Migration{}, fmt.Errorf("migration %s missing -- +goose Down marker", filename)
	}

	upSQL := strings.TrimSpace(parts[0])
	downSQL := strings.TrimSpace(parts[1])

	upSQL = strings.Replace(upSQL, "-- +goose Up", "", 1)
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
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_embedded BOOLEAN NOT NULL DEFAULT FALSE
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
	INSERT INTO schema_migrations (version, filename, is_embedded) 
	VALUES (?, ?, ?)`

	if _, err := tx.Exec(insertQuery, migration.Version, migration.Filename, migration.IsEmbedded); err != nil {
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
