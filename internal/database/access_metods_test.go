package database_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/GoldenDeals/DepGit/internal/config"
	. "github.com/GoldenDeals/DepGit/internal/database"
	dberror "github.com/GoldenDeals/DepGit/internal/share/errors"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/google/uuid"
	ase "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var log = logger.New("db_tests")

// setupTestDB creates a new database for an individual test
func setupTestDB(t *testing.T) (*DB, func()) {
	// Create a temp dir for the test
	tmpDir, err := os.MkdirTemp("", "depgit-test-*")
	require.NoError(t, err)

	// Create migrations directory
	migrationsDir := filepath.Join(tmpDir, "migrations")
	err = os.Mkdir(migrationsDir, 0755)
	require.NoError(t, err)

	// Create initial schema migration file
	initialMigrationFile := filepath.Join(migrationsDir, "001_initial_schema.sql")
	initialSchema := `
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created DATETIME NOT NULL,
    edited DATETIME,
    deleted DATETIME
);

CREATE TABLE IF NOT EXISTS keys (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type INTEGER NOT NULL,
    data BLOB NOT NULL,
    created DATETIME NOT NULL,
    deleted DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS permitions (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created DATETIME NOT NULL,
    edited DATETIME,
    deleted DATETIME
);

CREATE TABLE IF NOT EXISTS roles (
    role_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    rep_id TEXT NOT NULL,
    branch TEXT,
    created DATETIME NOT NULL,
    deleted DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (rep_id) REFERENCES permitions(id)
);`
	err = os.WriteFile(initialMigrationFile, []byte(initialSchema), 0644)
	require.NoError(t, err)

	// Create a DB in the temp dir
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create test configuration
	cfg := &config.Configuration{
		DB: config.DBConfig{
			Path:             dbPath,
			MigrationsPath:   migrationsDir,
			InitialMigration: initialSchema,
		},
	}

	// Initialize the database
	db := &DB{}
	err = db.Init(cfg)
	require.NoError(t, err)

	// Return db and cleanup function
	return db, func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}
}

func TestCreateUser(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	user := NewUser("Jonson", "gayporno@yandex.ru")
	err := db.CreateUser(ctx, &user)
	assert.Nil(err)

	checkuser, err := db.GetUser(ctx, user.ID)
	assert.Nil(err)

	// Compare only the fields we care about
	assert.Equal(user.ID, checkuser.ID)
	assert.Equal(user.Name, checkuser.Name)
	assert.Equal(user.Email, checkuser.Email)
}

func TestEditUser(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	user := NewUser("Jonson", "gayporno@yandex.ru")
	checkuser := NewUser("Soprano", "gayporno@gmail.ru")

	err := db.CreateUser(ctx, &user)
	assert.Nil(err)

	err = db.EditUser(ctx, user.ID, &checkuser)
	assert.Nil(err)

	checkwithuser, err := db.GetUser(ctx, user.ID)
	assert.Nil(err)

	// Compare only the fields we care about
	assert.Equal(checkuser.Name, checkwithuser.Name)
	assert.Equal(checkuser.Email, checkwithuser.Email)
}

func TestDeleteUser(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	user := NewUser("Jonson", "gayporno@yandex.ru")

	err := db.CreateUser(ctx, &user)
	assert.Nil(err)

	err = db.DeleteUser(ctx, user.ID)
	assert.Nil(err)

	// Verify the user is deleted
	_, err = db.GetUser(ctx, user.ID)
	assert.Equal(dberror.ErrNotFound, err)
}

func TestGetUser(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	user := NewUser("Jonson", "gayporno@yandex.ru")

	err := db.CreateUser(ctx, &user)
	assert.Nil(err)

	checkuser, err := db.GetUser(ctx, user.ID)
	assert.Nil(err)

	// Compare only the fields we care about
	assert.Equal(user.ID, checkuser.ID)
	assert.Equal(user.Name, checkuser.Name)
	assert.Equal(user.Email, checkuser.Email)
}

func TestGetUsers(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create multiple users
	user1 := NewUser("User1", "user1@example.com")
	user2 := NewUser("User2", "user2@example.com")
	user3 := NewUser("User3", "user3@example.com")

	err := db.CreateUser(ctx, &user1)
	assert.Nil(err)
	err = db.CreateUser(ctx, &user2)
	assert.Nil(err)
	err = db.CreateUser(ctx, &user3)
	assert.Nil(err)

	// Get all users
	users, err := db.GetUsers(ctx)
	assert.Nil(err)

	// We should have exactly 3 users
	assert.Equal(3, len(users))

	// Check if our created users are in the list
	found1, found2, found3 := false, false, false
	for _, u := range users {
		if u.ID == user1.ID {
			found1 = true
			assert.Equal(user1.Name, u.Name)
			assert.Equal(user1.Email, u.Email)
		}
		if u.ID == user2.ID {
			found2 = true
			assert.Equal(user2.Name, u.Name)
			assert.Equal(user2.Email, u.Email)
		}
		if u.ID == user3.ID {
			found3 = true
			assert.Equal(user3.Name, u.Name)
			assert.Equal(user3.Email, u.Email)
		}
	}

	assert.True(found1, "User1 should be in the list")
	assert.True(found2, "User2 should be in the list")
	assert.True(found3, "User3 should be in the list")
}

func TestAddSshKey(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	user := NewUser("Ahmed", "romalox@yandex.ru")
	data := make([]byte, 16)
	key := NewSShKey("lololoshka", RSA_SHA2_256, data)

	err := db.CreateUser(ctx, &user)
	assert.Nil(err)

	err = db.AddSshKey(ctx, user.ID, &key)
	assert.Nil(err)
}

func TestDeleteSshKey(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	user := NewUser("Eblan", "bibobo@gmail.ru")
	data := make([]byte, 16)
	key := NewSShKey("lololoshka", RSA_SHA2_256, data)
	err := db.CreateUser(ctx, &user)
	assert.Nil(err)

	err = db.AddSshKey(ctx, user.ID, &key)
	assert.Nil(err)

	err = db.DeleteSshKey(ctx, key.ID)
	assert.Nil(err)
}

func TestGetSshKeys(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user with multiple SSH keys
	user := NewUser("KeyUser", "keyuser@example.com")
	err := db.CreateUser(ctx, &user)
	assert.Nil(err)

	// Create several SSH keys
	data1 := []byte("ssh-key-data-1")
	data2 := []byte("ssh-key-data-2")
	data3 := []byte("ssh-key-data-3")

	key1 := NewSShKey("key1", RSA_SHA2_256, data1)
	key2 := NewSShKey("key2", SSH_RSA, data2)
	key3 := NewSShKey("key3", ECDSA_SHA2_NISTP256, data3)

	// Add all keys to the user
	err = db.AddSshKey(ctx, user.ID, &key1)
	assert.Nil(err)
	err = db.AddSshKey(ctx, user.ID, &key2)
	assert.Nil(err)
	err = db.AddSshKey(ctx, user.ID, &key3)
	assert.Nil(err)

	// Get all keys for the user
	keys, err := db.GetSshKeys(ctx, user.ID)
	assert.Nil(err)

	// We should have 3 keys
	assert.Equal(3, len(keys))

	// Verify each key exists in the result
	foundKeys := map[string]bool{
		key1.ID.String(): false,
		key2.ID.String(): false,
		key3.ID.String(): false,
	}

	for _, k := range keys {
		foundKeys[k.ID.String()] = true

		// Check the key properties
		if k.ID == key1.ID {
			assert.Equal(key1.Name, k.Name)
			assert.Equal(key1.Type, k.Type)
			assert.Equal(key1.Data, k.Data)
		} else if k.ID == key2.ID {
			assert.Equal(key2.Name, k.Name)
			assert.Equal(key2.Type, k.Type)
			assert.Equal(key2.Data, k.Data)
		} else if k.ID == key3.ID {
			assert.Equal(key3.Name, k.Name)
			assert.Equal(key3.Type, k.Type)
			assert.Equal(key3.Data, k.Data)
		}
	}

	for keyID, found := range foundKeys {
		assert.True(found, "Key %s should be found", keyID)
	}
}

func TestRepoOperations(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a new repository
	repo := NewRepo("test-repo")
	err := db.CreateRepo(ctx, &repo)
	assert.Nil(err)

	// Verify the repo has an ID
	assert.NotEqual(uuid.Nil, repo.ID)
}

func TestUpdateRepo(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a new repository
	repo := NewRepo("original-repo-name")
	err := db.CreateRepo(ctx, &repo)
	assert.Nil(err)

	// Update the repository
	repo.Name = "updated-repo-name"
	err = db.UpdateRepo(ctx, repo.ID, repo)
	assert.Nil(err)

	// Get the repository and verify the changes
	updatedRepo, err := db.GetRepo(ctx, repo.ID)
	assert.Nil(err)
	assert.Equal("updated-repo-name", updatedRepo.Name)
}

func TestDeleteRepo(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a new repository
	repo := NewRepo("repo-to-delete")
	err := db.CreateRepo(ctx, &repo)
	assert.Nil(err)

	// Delete the repository
	err = db.DeleteRepo(ctx, repo.ID)
	assert.Nil(err)

	// Try to get the repo (should fail)
	_, err = db.GetRepo(ctx, repo.ID)
	assert.NotNil(err)
}

func TestGetRepo(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a new repository
	repo := NewRepo("specific-repo")
	err := db.CreateRepo(ctx, &repo)
	assert.Nil(err)

	// Get the repository
	fetchedRepo, err := db.GetRepo(ctx, repo.ID)
	assert.Nil(err)

	// Compare the repositories
	assert.Equal(repo.ID, fetchedRepo.ID)
	assert.Equal(repo.Name, fetchedRepo.Name)
}

func TestGetRepos(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create multiple repositories
	repo1 := NewRepo("repo1")
	repo2 := NewRepo("repo2")
	repo3 := NewRepo("repo3")

	err1 := db.CreateRepo(ctx, &repo1)
	assert.Nil(err1)
	err2 := db.CreateRepo(ctx, &repo2)
	assert.Nil(err2)
	err3 := db.CreateRepo(ctx, &repo3)
	assert.Nil(err3)

	// Get all repositories
	repos, err := db.GetRepos(ctx)
	assert.Nil(err)

	// We should have exactly 3 repositories
	assert.Equal(3, len(repos))

	// Check if our created repos are in the list
	found1, found2, found3 := false, false, false
	for _, r := range repos {
		if r.ID == repo1.ID {
			found1 = true
			assert.Equal(repo1.Name, r.Name)
		}
		if r.ID == repo2.ID {
			found2 = true
			assert.Equal(repo2.Name, r.Name)
		}
		if r.ID == repo3.ID {
			found3 = true
			assert.Equal(repo3.Name, r.Name)
		}
	}

	assert.True(found1, "Repo1 should be in the list")
	assert.True(found2, "Repo2 should be in the list")
	assert.True(found3, "Repo3 should be in the list")
}

func TestAccessRoleCRUD(t *testing.T) {
	assert := ase.New(t)
	require := require.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user and a repo for testing
	user := NewUser("RoleUser", "roleuser@example.com")
	err := db.CreateUser(ctx, &user)
	assert.Nil(err)

	repo := NewRepo("role-test-repo")
	err = db.CreateRepo(ctx, &repo)
	assert.Nil(err)

	// Create an access role
	role := AccessRole{
		RoleID: uuid.New(),
		UserID: user.ID,
		RepoID: repo.ID,
	}

	// Test CreateAccessRole
	err = db.CreateAccessRole(ctx, &role)
	require.Nil(err)

	// Test GetAccessRole
	fetchedRole, err := db.GetAccessRole(ctx, role.RoleID)
	require.Nil(err)
	assert.Equal(role.RoleID, fetchedRole.RoleID)
	assert.Equal(role.UserID, fetchedRole.UserID)
	assert.Equal(role.RepoID, fetchedRole.RepoID)

	// Test DeleteAccessRole
	err = db.DeleteAccessRole(ctx, role.RoleID)
	assert.Nil(err)

	// Verify role is deleted
	_, err = db.GetAccessRole(ctx, role.RoleID)
	assert.NotNil(err)
}

func TestCheckPermissions(t *testing.T) {
	assert := ase.New(t)
	ctx := context.Background()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user and a repo for testing
	user := NewUser("PermUser", "permuser@example.com")
	err := db.CreateUser(ctx, &user)
	assert.Nil(err)

	repo := NewRepo("perm-test-repo")
	err = db.CreateRepo(ctx, &repo)
	assert.Nil(err)

	// Create an access role with a specific branch
	role := AccessRole{
		RoleID:   uuid.New(),
		UserID:   user.ID,
		RepoID:   repo.ID,
		Branches: nil, // No specific branch restrictions
		Created:  time.Now(),
	}

	err = db.CreateAccessRole(ctx, &role)
	assert.Nil(err)

	// Check permissions without specifying branch
	hasPermission, err := db.CheckPermissions(ctx, user.ID, repo.ID, "")
	assert.Nil(err)
	assert.True(hasPermission, "User should have permission for any branch")

	// Check with non-existent user
	hasPermission, err = db.CheckPermissions(ctx, uuid.New(), repo.ID, "")
	assert.Nil(err)
	assert.False(hasPermission, "Non-existent user should not have permission")
}
