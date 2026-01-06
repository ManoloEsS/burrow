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

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("could not set goose dialect: %w", err)
	}
	if err := goose.Up(db, dbMigrations); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

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

func (db *Database) WithContext(ctx context.Context, timeoutSec time.Duration) (context.Context, context.CancelFunc) {
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	return context.WithTimeout(ctx, timeoutSec)
}
