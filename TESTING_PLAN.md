# Testing Expansion Readme

This document tracks the planned tests that will be added before implementation begins. Each entry includes the tentative name and the specific behavior it should verify. We will update this file as plans evolve and mark items complete once the corresponding tests land.

## Unit Tests

1. ~~**config.Load defaults** – Exercise `Load` when the config file is missing or malformed, confirm environment variable overrides take precedence, and assert `EnsureDirectories` surfaces OS errors.~~ ✅ *Implemented: `config_test.go` now isolates XDG dirs, adds YAML helpers, verifies env overrides, and forces `mkdirAll` failures.*
2. ~~**config.Paths edge cases** – Cover `SearchConfigFile`, `ConfigFileExists`, and `EnsureDirectories` when XDG directories are absent or unwritable to guarantee correct fallback behavior.~~ ✅ *Implemented: `paths_test.go` now stubs XDG roots, checks config discovery, existence detection, and directory creation.*
3. **domain.Response failures** – Validate `BuildResponse` when body reads fail, JSON pretty-printing errors occur, or large non-JSON payloads are returned to ensure graceful degradation.
4. **domain.Request validations** – Add negative cases for `ParseMethod`, `ParseUrl`, `ParseHeaders`, and `ParseBody` to ensure bad inputs are rejected with actionable errors.
5. **service.HttpClient SendRequest** – Mock the HTTP transport to verify timeout handling, automatic header defaults, query parameter merging, and accurate response-time recording.
6. **service.HttpClient persistence CRUD** – Use an in-memory SQLite database to test `SaveRequest`, `GetSavedRequests`, and `DeleteRequest`, covering serialization errors and deletion of missing entries.
7. **service.Server cleanupBinary** – Ensure `cleanupBinary` removes existing binaries, reports filesystem errors, and resets internal state.
8. **service.Server healthChecker failures** – Inject a stubbed HTTP client to simulate non-200 responses and timeouts, verifying retries and emitted UI events.
9. **database.MigrationRunner core** – Unit test `parseMigration`, `LoadMigrations`, and `Up` with an in-memory database to confirm idempotence and proper error propagation.

## Integration Tests

1. **service.Server lifecycle** – Compile and run a temporary Go server file, then use `StartServer`/`StopServer` to assert build → start → health check → shutdown flow, including emitted events and cleanup.
2. **database.MigrationRunner end-to-end** – Run embedded migrations against a temporary SQLite file, verify schema state, and ensure rerunning migrations is safe.
3. **cmd/burrow main flow** – Execute the CLI against temporary config/database paths to ensure flag parsing, config loading, migration execution, and graceful termination all function together.

## UI / Other Tests

1. **tui event handling** – (If applicable) simulate UI components reacting to `ServerService` events to guard state transitions and rendering logic from regressions.

## Workflow Reminder

We will plan each test in detail before implementation. No code changes for these tests will occur until explicit approval is given, and every new test will be followed by running the relevant test suite to verify correctness.
