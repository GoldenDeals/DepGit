package migrations_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoldenDeals/DepGit/internal/config"
	"github.com/GoldenDeals/DepGit/internal/database"
	"github.com/GoldenDeals/DepGit/internal/database/migrations"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationsIntegration(t *testing.T) {
	// Create a temporary directory for test database and migrations
	tempDir, err := os.MkdirTemp("", "migrations-integration-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create migrations directory
	migrationsDir := filepath.Join(tempDir, "migrations")
	err = os.Mkdir(migrationsDir, 0755)
	require.NoError(t, err)

	// Create a test migration file
	migrationFile := filepath.Join(migrationsDir, "001_create_test_table.sql")
	err = os.WriteFile(migrationFile, []byte(`
		CREATE TABLE test_table (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		);
	`), 0644)
	require.NoError(t, err)

	// Create a second migration file
	secondMigrationFile := filepath.Join(migrationsDir, "002_add_index.sql")
	err = os.WriteFile(secondMigrationFile, []byte(`
		CREATE INDEX idx_test_table_name ON test_table(name);
	`), 0644)
	require.NoError(t, err)

	// Create the database path
	dbPath := filepath.Join(tempDir, "test.db")

	// Create test configuration
	cfg := &config.Configuration{
		DB: config.DBConfig{
			Path:             dbPath,
			MigrationsPath:   migrationsDir,
			InitialMigration: "",
		},
	}

	// Initialize the database
	db := &database.DB{}
	err = db.Init(cfg)
	require.NoError(t, err)
	defer db.Close()

	// Verify the test_table was created
	sqlDB, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer sqlDB.Close()

	var tableCount int
	err = sqlDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&tableCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, tableCount)

	// Verify the index was created
	var indexCount int
	err = sqlDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='index' AND name='idx_test_table_name'").Scan(&indexCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, indexCount)

	// Add a third migration directly using the migrations API
	thirdMigrationFile := filepath.Join(migrationsDir, "003_add_column.sql")
	err = os.WriteFile(thirdMigrationFile, []byte(`
		ALTER TABLE test_table ADD COLUMN description TEXT;
	`), 0644)
	require.NoError(t, err)

	// Run the migrations manually
	migrationManager := migrations.New(sqlDB)
	err = migrationManager.Initialize(context.Background())
	require.NoError(t, err)

	err = migrationManager.LoadFromDir(context.Background(), migrationsDir)
	require.NoError(t, err)

	err = migrationManager.Apply(context.Background())
	require.NoError(t, err)

	// Verify the column was added
	var columnCount int
	err = sqlDB.QueryRow("SELECT count(*) FROM pragma_table_info('test_table') WHERE name='description'").Scan(&columnCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, columnCount)

	// Verify all migrations were recorded
	var migrationCount int
	err = sqlDB.QueryRow("SELECT count(*) FROM migrations").Scan(&migrationCount)
	assert.NoError(t, err)
	assert.Equal(t, 3, migrationCount)
}
