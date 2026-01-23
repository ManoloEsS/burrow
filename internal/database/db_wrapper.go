package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB      *sql.DB
	Queries *Queries
	Timeout time.Duration
}

func NewDatabase(dbPath, dbString string) (*Database, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("could not create database directory: %w", err)
	}

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

	migrationRunner := NewMigrationRunner(db, "")
	if err := migrationRunner.LoadMigrations(); err != nil {
		return nil, fmt.Errorf("failed to load migrations: %w", err)
	}
	if err := migrationRunner.Up(); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	queries := New(db)

	log.Printf("Database initialized at %s", dbPath)
	return &Database{
		DB:      db,
		Queries: queries,
		Timeout: 5 * time.Second,
	}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}

func (db *Database) WithContext(ctx context.Context, timeoutSec time.Duration) (context.Context, context.CancelFunc) {
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	return context.WithTimeout(ctx, timeoutSec)
}
