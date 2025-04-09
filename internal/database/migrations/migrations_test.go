package migrations

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDb(t *testing.T) (*sql.DB, func()) {
	// Create a temporary file for the test database
	tmpFile, err := os.CreateTemp("", "test-migrations-*.db")
	require.NoError(t, err)
	tmpFile.Close()

	// Open the database
	db, err := sql.Open("sqlite3", tmpFile.Name())
	require.NoError(t, err)

	// Return the database and a cleanup function
	return db, func() {
		db.Close()
		os.Remove(tmpFile.Name())
	}
}

func createTestMigrationDir(t *testing.T) (string, func()) {
	// Create a temporary directory for test migrations
	tmpDir, err := os.MkdirTemp("", "test-migrations-dir")
	require.NoError(t, err)

	// Create a test migration file
	migration1 := filepath.Join(tmpDir, "001_test_migration.sql")
	err = os.WriteFile(migration1, []byte(`
		CREATE TABLE test_table (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);
	`), 0644)
	require.NoError(t, err)

	// Create a second test migration file
	migration2 := filepath.Join(tmpDir, "002_test_migration.sql")
	err = os.WriteFile(migration2, []byte(`
		CREATE TABLE another_table (
			id INTEGER PRIMARY KEY,
			value TEXT NOT NULL
		);
	`), 0644)
	require.NoError(t, err)

	return tmpDir, func() {
		os.RemoveAll(tmpDir)
	}
}

func TestMigrationManager_Initialize(t *testing.T) {
	db, cleanup := setupTestDb(t)
	defer cleanup()

	// Create a migration manager
	manager := New(db)

	// Initialize the migrations table
	ctx := context.Background()
	err := manager.Initialize(ctx)
	assert.NoError(t, err)

	// Verify the migrations table was created
	var count int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='migrations'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestMigrationManager_LoadFromDir(t *testing.T) {
	db, dbCleanup := setupTestDb(t)
	defer dbCleanup()

	migrationDir, dirCleanup := createTestMigrationDir(t)
	defer dirCleanup()

	// Create a migration manager
	manager := New(db)

	// Load migrations from directory
	ctx := context.Background()
	err := manager.LoadFromDir(ctx, migrationDir)
	assert.NoError(t, err)

	// Verify migrations were loaded
	migrations := manager.GetMigrations()
	assert.Len(t, migrations, 2)
	assert.Equal(t, "001_test_migration.sql", migrations[0].Name)
	assert.Equal(t, "002_test_migration.sql", migrations[1].Name)
}

func TestMigrationManager_Apply(t *testing.T) {
	db, dbCleanup := setupTestDb(t)
	defer dbCleanup()

	migrationDir, dirCleanup := createTestMigrationDir(t)
	defer dirCleanup()

	// Create a migration manager
	manager := New(db)

	// Initialize the migrations table
	ctx := context.Background()
	err := manager.Initialize(ctx)
	assert.NoError(t, err)

	// Load migrations from directory
	err = manager.LoadFromDir(ctx, migrationDir)
	assert.NoError(t, err)

	// Apply migrations
	err = manager.Apply(ctx)
	assert.NoError(t, err)

	// Verify the test tables were created
	var count int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='another_table'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Verify migrations were recorded
	var migrationCount int
	err = db.QueryRow("SELECT count(*) FROM migrations").Scan(&migrationCount)
	assert.NoError(t, err)
	assert.Equal(t, 2, migrationCount)

	// Apply migrations again should be a no-op
	err = manager.Apply(ctx)
	assert.NoError(t, err)

	// Migration count should still be 2
	err = db.QueryRow("SELECT count(*) FROM migrations").Scan(&migrationCount)
	assert.NoError(t, err)
	assert.Equal(t, 2, migrationCount)
}

func TestMigrationManager_AddMigration(t *testing.T) {
	db, cleanup := setupTestDb(t)
	defer cleanup()

	// Create a migration manager
	manager := New(db)

	// Add a migration
	manager.AddMigration("003_manual_migration.sql", "CREATE TABLE manual_table (id INTEGER PRIMARY KEY);")

	// Verify the migration was added
	migrations := manager.GetMigrations()
	assert.Len(t, migrations, 1)
	assert.Equal(t, "003_manual_migration.sql", migrations[0].Name)

	// Initialize and apply the migration
	ctx := context.Background()
	err := manager.Initialize(ctx)
	assert.NoError(t, err)

	err = manager.Apply(ctx)
	assert.NoError(t, err)

	// Verify the table was created
	var count int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='manual_table'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
