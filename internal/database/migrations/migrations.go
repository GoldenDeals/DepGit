package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/GoldenDeals/DepGit/internal/share/logger"
)

// Logger for migrations
var migrationLogger = logger.New("db_migrations")

// Migration represents a database migration
type Migration struct {
	ID        int64
	Name      string
	SQL       string
	AppliedAt time.Time
}

// Manager handles database migrations
type Manager struct {
	db         *sql.DB
	migrations []*Migration
}

// New creates a new migration manager
func New(db *sql.DB) *Manager {
	return &Manager{
		db:         db,
		migrations: []*Migration{},
	}
}

// Initialize creates the migrations table if it doesn't exist
func (m *Manager) Initialize(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	_, err := m.db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}

// LoadFromDir loads migrations from a directory
func (m *Manager) LoadFromDir(ctx context.Context, dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrationLogger.
		WithContext(ctx).
		WithField("directory", dir).
		Info("Loading migrations from directory")

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		path := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		m.migrations = append(m.migrations, &Migration{
			Name: file.Name(),
			SQL:  string(content),
		})
	}

	// Sort migrations by name
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Name < m.migrations[j].Name
	})

	return nil
}

// GetAppliedMigrations returns a list of applied migrations
func (m *Manager) GetAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT name FROM migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		applied[name] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating migrations: %w", err)
	}

	return applied, nil
}

// Apply applies all pending migrations
func (m *Manager) Apply(ctx context.Context) error {
	applied, err := m.GetAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			migrationLogger.WithContext(ctx).WithError(rbErr).Error("Failed to rollback transaction")
		}
	}()

	for _, migration := range m.migrations {
		if applied[migration.Name] {
			migrationLogger.
				WithContext(ctx).
				WithField("migration", migration.Name).
				Info("Migration already applied, skipping")
			continue
		}

		migrationLogger.
			WithContext(ctx).
			WithField("migration", migration.Name).
			Info("Applying migration")

		// Split the SQL into individual statements
		statements := splitStatements(migration.SQL)

		for _, stmt := range statements {
			if strings.TrimSpace(stmt) == "" {
				continue
			}

			_, err = tx.ExecContext(ctx, stmt)
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", migration.Name, err)
			}
		}

		// Record the migration
		_, err = tx.ExecContext(ctx, "INSERT INTO migrations (name) VALUES (?)", migration.Name)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
		}

		migrationLogger.
			WithContext(ctx).
			WithField("migration", migration.Name).
			Info("Migration applied successfully")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// AddMigration adds a migration to the manager
func (m *Manager) AddMigration(name, sql string) {
	m.migrations = append(m.migrations, &Migration{
		Name: name,
		SQL:  sql,
	})
}

// GetMigrations returns the list of migrations
func (m *Manager) GetMigrations() []*Migration {
	return m.migrations
}

// splitStatements splits a SQL script into individual statements
func splitStatements(script string) []string {
	// Simple implementation - split by semicolon
	// This doesn't handle all edge cases (like semicolons in strings)
	// but works for most simple migration scripts
	statements := strings.Split(script, ";")

	// Clean up statements
	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}
