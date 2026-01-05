package store

import (
	"database/sql"
	"fmt"
)

type Migration struct {
	Version int
	Up      string
	Down    string
}

var migrations = []Migration{
	{
		Version: 1,
		Up: `
			CREATE TABLE IF NOT EXISTS processed_results (
				id TEXT PRIMARY KEY,
				request_json JSONB NOT NULL,
				response_json JSONB NOT NULL,
				topic TEXT NOT NULL,
				topic_source TEXT NOT NULL,
				topic_confidence REAL NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE INDEX IF NOT EXISTS idx_processed_results_created_at ON processed_results(created_at);
			CREATE INDEX IF NOT EXISTS idx_processed_results_topic ON processed_results(topic);
		`,
		Down: `
			DROP INDEX IF EXISTS idx_processed_results_topic;
			DROP INDEX IF EXISTS idx_processed_results_created_at;
			DROP TABLE IF EXISTS processed_results;
		`,
	},
}

func runMigrations(db *sql.DB) error {
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for _, migration := range migrations {
		applied, err := isMigrationApplied(db, migration.Version)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if applied {
			continue
		}

		if err := applyMigration(db, migration); err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	_, err := db.Exec(query)
	return err
}

func isMigrationApplied(db *sql.DB, version int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func applyMigration(db *sql.DB, migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(migration.Up); err != nil {
		return fmt.Errorf("failed to execute migration up: %w", err)
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

