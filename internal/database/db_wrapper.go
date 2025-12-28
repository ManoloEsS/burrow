package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Database struct {
	*sql.DB
	Queries *Queries
	Timeout time.Duration
}

// NewDatabase opens the DB, applies Goose migrations (Go mode), and dumps schema for SQLC
func NewDatabase(dbPath, dbString, dbMigrations string) (*Database, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("could not create database directory: %w", err)
	}

	// Open DB with pragmas in DSN
	db, err := sql.Open("sqlite3", dbString)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	// Apply migrations via Go mode
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("could not set goose dialect: %w", err)
	}
	if err := goose.Up(db, dbMigrations); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Dump schema for SQLC
	if err := dumpSchema(db, "sql/schema/schema.sql"); err != nil {
		return nil, fmt.Errorf("failed to dump schema: %w", err)
	}

	queries := New(db)

	log.Printf("Database initialized at %s", dbPath)
	return &Database{
		DB:      db,
		Queries: queries,
		Timeout: 5 * time.Second,
	}, nil
}

// dumpSchema writes current DB schema to file (for SQLC)
func dumpSchema(db *sql.DB, path string) error {
	rows, err := db.Query(`SELECT sql 
FROM sqlite_master 
WHERE type IN ('table','index','trigger')
AND name NOT IN ('goose_db_version','sqlite_sequence');`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var schema []string
	for rows.Next() {
		var sqlStmt sql.NullString
		if err := rows.Scan(&sqlStmt); err != nil {
			return err
		}
		if sqlStmt.Valid && strings.TrimSpace(sqlStmt.String) != "" {
			schema = append(schema, sqlStmt.String+";\n")
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(strings.Join(schema, "\n")), 0o644)
}

// Close safely closes DB
func (db *Database) Close() error {
	return db.DB.Close()
}

// WithContext returns a context with timeout
func (db *Database) WithContext(ctx context.Context, timeoutSec time.Duration) (context.Context, context.CancelFunc) {
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	return context.WithTimeout(ctx, timeoutSec)
}
