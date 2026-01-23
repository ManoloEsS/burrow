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
)

type Database struct {
	DB      *sql.DB
	Queries *Queries
	Timeout time.Duration
}

func NewDatabase(dbPath, dbString, dbMigrations string) (*Database, error) {
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

	migrationRunner := NewMigrationRunner(db, dbMigrations)
	if err := migrationRunner.LoadMigrations(); err != nil {
		return nil, fmt.Errorf("failed to load migrations: %w", err)
	}
	if err := migrationRunner.Up(); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	if isRunningFromSource() {
		if err := dumpSchema(db, "sql/schema/schema.sql"); err != nil {
			log.Printf("Warning: failed to dump schema: %v", err)
		}
	}

	queries := New(db)

	log.Printf("Database initialized at %s", dbPath)
	return &Database{
		DB:      db,
		Queries: queries,
		Timeout: 5 * time.Second,
	}, nil
}

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

func (db *Database) Close() error {
	return db.DB.Close()
}

func isRunningFromSource() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}

	exeDir := filepath.Dir(exePath)

	goModPath := filepath.Join(exeDir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return true
	}

	goModPath = filepath.Join(filepath.Dir(exeDir), "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return true
	}

	return false
}

func (db *Database) WithContext(ctx context.Context, timeoutSec time.Duration) (context.Context, context.CancelFunc) {
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	return context.WithTimeout(ctx, timeoutSec)
}
