.PHONY: build run generate clean

# Generate SQLC code
generate:
	sqlc generate

# Run application
run:
	go run cmd/cli/main.go

# Build application
build:
	go build -o burrow cmd/cli/main.go

# Clean artifacts
clean:
	rm -f burrow *.db sql/schema.sql

# Full setup
setup: clean run generate

# Check migration status
migrate-status:
	@echo "Checking migration status..."
	sqlite3 ./burrow.db "SELECT version FROM goose_db_version;"

# List SQLC generated files
list-generated:
	@echo "Generated files:"
	ls -la internal/database/