package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoldenDeals/DepGit/internal/database"
	"github.com/GoldenDeals/DepGit/internal/gen/api"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	testDB     *database.DB
	testDBPath string
)

// TestMain is the entry point for all tests in this package
func TestMain(m *testing.M) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		// If .env file doesn't exist, create a default one for testing
		log.Println("Warning: .env file not found, using default test settings")
		testDBPath = "./test_depgit.db"
	} else {
		// Get test database path from environment variable
		testDBPath = os.Getenv("TEST_DB_PATH")
		if testDBPath == "" {
			testDBPath = "./test_depgit.db" // Default if not specified
		}
	}

	// Setup test database
	testDB = setupTestDB()

	// Run tests
	code := m.Run()

	// Cleanup
	teardownTestDB()

	// Exit with test result code
	os.Exit(code)
}

func setupTestDB() *database.DB {
	// Remove any existing test database
	os.Remove(testDBPath)

	// Create a new database
	db := &database.DB{}
	err := db.Init(testDBPath)
	if err != nil {
		panic(err)
	}

	// Apply migrations from main.sql file
	migrationSQL, err := os.ReadFile("../../migrations/main.sql")
	if err != nil {
		// If main.sql doesn't exist, create tables manually
		log.Println("Warning: migrations/main.sql not found, creating tables manually")
		_, err = db.DB().Exec(`
			CREATE TABLE users (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				email TEXT NOT NULL,
				created_at DATETIME,
				updated_at DATETIME,
				deleted_at DATETIME
			);

			CREATE TABLE keys (
				id TEXT PRIMARY KEY,
				user_id TEXT NOT NULL,
				name TEXT NOT NULL,
				type INTEGER NOT NULL,
				key BLOB NOT NULL,
				created_at DATETIME,
				deleted_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id)
			);
		`)
		if err != nil {
			panic(fmt.Sprintf("Failed to create tables: %v", err))
		}
	} else {
		// Execute the migration SQL
		_, err = db.DB().Exec(string(migrationSQL))
		if err != nil {
			panic(fmt.Sprintf("Failed to apply migrations: %v", err))
		}
	}

	// Store the database connection for tests to use
	testDB = db

	absPath, _ := filepath.Abs(testDBPath)
	log.Printf("Test database initialized at %s", absPath)

	return db
}

func teardownTestDB() {
	// Close the database connection
	if testDB != nil {
		testDB.Close()
	}

	// Remove the database file
	os.Remove(testDBPath)
	log.Printf("Test database removed: %s", testDBPath)
}

func TestGetUsers(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new API handler with the test database
	handler := &APIHandler{db: testDB}

	// Create a new request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := handler.GetUsers(c, api.GetUsersParams{})

	// Assert no error
	assert.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var users []api.User
	err = json.Unmarshal(rec.Body.Bytes(), &users)
	assert.NoError(t, err)

	// We don't know how many users there are, but the response should be a valid JSON array
	assert.IsType(t, []api.User{}, users)
}

func TestCreateUser(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new API handler with the test database
	handler := &APIHandler{db: testDB}

	// Create a new user
	user := api.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Convert user to JSON
	userJSON, err := json.Marshal(user)
	assert.NoError(t, err)

	// Create a new request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err = handler.CreateUser(c)

	// Assert no error
	assert.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Parse response
	var createdUser api.User
	err = json.Unmarshal(rec.Body.Bytes(), &createdUser)
	assert.NoError(t, err)

	// Assert user properties
	assert.Equal(t, user.Username, createdUser.Username)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.NotNil(t, createdUser.Id)
	assert.NotNil(t, createdUser.Role)
	assert.NotNil(t, createdUser.CreatedAt)
}

func TestGetUser(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new API handler with the test database
	handler := &APIHandler{db: testDB}

	// First, create a user to get
	user := api.User{
		Username: "testuser2",
		Email:    "test2@example.com",
	}

	// Convert user to JSON
	userJSON, err := json.Marshal(user)
	assert.NoError(t, err)

	// Create a new request to create the user
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler to create the user
	err = handler.CreateUser(c)
	assert.NoError(t, err)

	// Parse response to get the created user's ID
	var createdUser api.User
	err = json.Unmarshal(rec.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser.Id)

	// Now, get the user by ID
	userIdStr := createdUser.Id.String()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userIdStr, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues(userIdStr)

	// Call the handler to get the user
	err = handler.GetUser(c, *createdUser.Id)

	// Assert no error
	assert.NoError(t, err)

	// Assert status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var retrievedUser api.User
	err = json.Unmarshal(rec.Body.Bytes(), &retrievedUser)
	assert.NoError(t, err)

	// Assert user properties
	assert.Equal(t, createdUser.Username, retrievedUser.Username)
	assert.Equal(t, createdUser.Email, retrievedUser.Email)
	assert.Equal(t, createdUser.Id, retrievedUser.Id)
}
