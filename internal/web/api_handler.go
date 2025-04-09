// Package web provides the web server and API handlers for the DepGit application.
package web

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/GoldenDeals/DepGit/internal/database"
	"github.com/GoldenDeals/DepGit/internal/gen/api"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

var webLogger = logger.New("web")

// APIHandler implements the ServerInterface from the generated API code
type APIHandler struct {
	db *database.DB
}

// NewAPIHandler creates a new API handler with the given database connection
func NewAPIHandler(db *database.DB) *APIHandler {
	return &APIHandler{
		db: db,
	}
}

// GetUsers handles the GET /users endpoint
func (h *APIHandler) GetUsers(ctx echo.Context, params api.GetUsersParams) error {
	// Get users from database
	users, err := h.db.GetUsers(ctx.Request().Context())
	if err != nil {
		webLogger.WithError(err).Error("Failed to get users")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to get users"),
		})
	}

	// Convert database users to API users
	apiUsers := make([]api.User, 0, len(users))
	for _, user := range users {
		apiUser := dbUserToAPIUser(user)

		// Filter by role if specified
		if params.Role != nil {
			if apiUser.Role == nil || string(*apiUser.Role) != string(*params.Role) {
				continue
			}
		}

		apiUsers = append(apiUsers, apiUser)
	}

	return ctx.JSON(http.StatusOK, apiUsers)
}

// CreateUser handles the POST /users endpoint
func (h *APIHandler) CreateUser(ctx echo.Context) error {
	// Parse request body
	var reqUser api.User
	if err := ctx.Bind(&reqUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Code:    intPtr(http.StatusBadRequest),
			Message: strPtr("Invalid request body"),
		})
	}

	// Validate required fields
	if reqUser.Username == "" || reqUser.Email == "" {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Code:    intPtr(http.StatusBadRequest),
			Message: strPtr("Username and email are required"),
		})
	}

	// Create database user
	dbUser := database.NewUser(reqUser.Username, string(reqUser.Email))

	// Set role if provided
	if reqUser.Role != nil {
		// TODO: Store role in database
	}

	// Save user to database
	err := h.db.CreateUser(ctx.Request().Context(), &dbUser)
	if err != nil {
		webLogger.WithError(err).Error("Failed to create user")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to create user"),
		})
	}

	// Convert to API user and return
	apiUser := dbUserToAPIUser(dbUser)
	return ctx.JSON(http.StatusCreated, apiUser)
}

// GetUser handles the GET /users/{userId} endpoint
func (h *APIHandler) GetUser(ctx echo.Context, userId openapi_types.UUID) error {
	// UUID is already a valid UUID object, just use it directly
	dbID := userId

	// Get user from database
	user, err := h.db.GetUser(ctx.Request().Context(), dbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, api.Error{
				Code:    intPtr(http.StatusNotFound),
				Message: strPtr("User not found"),
			})
		}
		webLogger.WithError(err).Error("Failed to get user")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to get user"),
		})
	}

	// Convert to API user and return
	apiUser := dbUserToAPIUser(user)
	return ctx.JSON(http.StatusOK, apiUser)
}

// UpdateUser handles the PUT /users/{userId} endpoint
func (h *APIHandler) UpdateUser(ctx echo.Context, userId openapi_types.UUID) error {
	// Parse request body
	var reqUser api.User
	if err := ctx.Bind(&reqUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Code:    intPtr(http.StatusBadRequest),
			Message: strPtr("Invalid request body"),
		})
	}

	// UUID is already a valid UUID object, just use it directly
	dbID := userId

	// Get existing user
	existingUser, err := h.db.GetUser(ctx.Request().Context(), dbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, api.Error{
				Code:    intPtr(http.StatusNotFound),
				Message: strPtr("User not found"),
			})
		}
		webLogger.WithError(err).Error("Failed to get user for update")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to get user for update"),
		})
	}

	// Update user fields
	if reqUser.Username != "" {
		existingUser.Name = reqUser.Username
	}
	if reqUser.Email != "" {
		existingUser.Email = string(reqUser.Email)
	}

	// Update role if provided
	if reqUser.Role != nil {
		// TODO: Store role in database
	}

	// Set updated time
	existingUser.Edited = time.Now()

	// Save updated user
	err = h.db.EditUser(ctx.Request().Context(), dbID, &existingUser)
	if err != nil {
		webLogger.WithError(err).Error("Failed to update user")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to update user"),
		})
	}

	// Convert to API user and return
	apiUser := dbUserToAPIUser(existingUser)
	return ctx.JSON(http.StatusOK, apiUser)
}

// DeleteUser handles the DELETE /users/{userId} endpoint
func (h *APIHandler) DeleteUser(ctx echo.Context, userId openapi_types.UUID) error {
	// UUID is already a valid UUID object, just use it directly
	dbID := userId

	// Delete user from database
	err := h.db.DeleteUser(ctx.Request().Context(), dbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, api.Error{
				Code:    intPtr(http.StatusNotFound),
				Message: strPtr("User not found"),
			})
		}
		webLogger.WithError(err).Error("Failed to delete user")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to delete user"),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// GetSshKeys handles the GET /users/{userId}/ssh-keys endpoint
func (h *APIHandler) GetSshKeys(ctx echo.Context, userId openapi_types.UUID) error {
	// UUID is already a valid UUID object, just use it directly
	dbID := userId

	// Check if user exists
	_, err := h.db.GetUser(ctx.Request().Context(), dbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, api.Error{
				Code:    intPtr(http.StatusNotFound),
				Message: strPtr("User not found"),
			})
		}
		webLogger.WithError(err).Error("Failed to get user for SSH keys")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to get user for SSH keys"),
		})
	}

	// Get SSH keys from database
	// TODO: Implement GetSshKeys in database package
	keys := []database.SshKey{}

	// Convert database keys to API keys
	apiKeys := make([]api.SshKey, 0, len(keys))
	for _, key := range keys {
		apiKeys = append(apiKeys, dbSshKeyToAPISshKey(key))
	}

	return ctx.JSON(http.StatusOK, apiKeys)
}

// AddSshKey handles the POST /users/{userId}/ssh-keys endpoint
func (h *APIHandler) AddSshKey(ctx echo.Context, userId openapi_types.UUID) error {
	// Parse request body
	var reqKey api.SshKey
	if err := ctx.Bind(&reqKey); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Code:    intPtr(http.StatusBadRequest),
			Message: strPtr("Invalid request body"),
		})
	}

	// UUID is already a valid UUID object, just use it directly
	dbID := userId

	// Check if user exists
	_, err := h.db.GetUser(ctx.Request().Context(), dbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, api.Error{
				Code:    intPtr(http.StatusNotFound),
				Message: strPtr("User not found"),
			})
		}
		webLogger.WithError(err).Error("Failed to get user for adding SSH key")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to get user for adding SSH key"),
		})
	}

	// Create database SSH key
	dbKey := database.SshKey{
		ID:      uuid.New(),
		UserID:  dbID,
		Name:    reqKey.Name,
		Data:    []byte(reqKey.Key),
		Created: time.Now(),
	}

	// Set key type if provided
	if reqKey.Type != nil {
		dbKey.Type = database.SSH_KEY_TYPE(*reqKey.Type)
	} else {
		dbKey.Type = database.SSH_KEY_TYPE_RSA // Default type
	}

	// Save SSH key to database
	err2 := h.db.AddSshKey(ctx.Request().Context(), dbID, &dbKey)
	if err2 != nil {
		webLogger.WithError(err2).Error("Failed to add SSH key")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to add SSH key"),
		})
	}

	// Convert to API key and return
	apiKey := dbSshKeyToAPISshKey(dbKey)
	return ctx.JSON(http.StatusCreated, apiKey)
}

// DeleteSshKey handles the DELETE /ssh-keys/{keyId} endpoint
func (h *APIHandler) DeleteSshKey(ctx echo.Context, keyId openapi_types.UUID) error {
	// UUID is already a valid UUID object, just use it directly
	dbID := keyId

	// Delete SSH key from database
	err := h.db.DeleteSshKey(ctx.Request().Context(), dbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, api.Error{
				Code:    intPtr(http.StatusNotFound),
				Message: strPtr("SSH key not found"),
			})
		}
		webLogger.WithError(err).Error("Failed to delete SSH key")
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Code:    intPtr(http.StatusInternalServerError),
			Message: strPtr("Failed to delete SSH key"),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Helper functions for converting between database and API models

func dbUserToAPIUser(dbUser database.User) api.User {
	// UUID is already a valid UUID object
	id := dbUser.ID
	email := openapi_types.Email(dbUser.Email)

	// Default role to Developer
	role := api.UserRoleDeveloper

	return api.User{
		Id:        &id,
		Username:  dbUser.Name,
		Email:     email,
		Role:      &role,
		CreatedAt: &dbUser.Created,
		UpdatedAt: &dbUser.Edited,
	}
}

func dbSshKeyToAPISshKey(dbKey database.SshKey) api.SshKey {
	// UUIDs are already valid UUID objects
	id := dbKey.ID
	userId := dbKey.UserID
	keyType := api.SshKeyType(dbKey.Type)

	return api.SshKey{
		Id:        &id,
		UserId:    userId,
		Name:      dbKey.Name,
		Key:       string(dbKey.Data),
		Type:      &keyType,
		CreatedAt: &dbKey.Created,
	}
}

// Helper functions for creating pointers to primitives

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int32 {
	i32 := int32(i)
	return &i32
}

// Placeholder implementations for unimplemented API methods

func (h *APIHandler) DeleteAccessRole(ctx echo.Context, roleId openapi_types.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) UpdateAccessRole(ctx echo.Context, roleId openapi_types.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) ForgotPassword(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) Login(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) Logout(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) RefreshToken(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) Register(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) ResetPassword(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) ValidateToken(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) GetGitInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) GetRepos(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) CreateRepo(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) DeleteRepo(ctx echo.Context, repoId openapi_types.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) GetRepo(ctx echo.Context, repoId openapi_types.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) UpdateRepo(ctx echo.Context, repoId openapi_types.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) GetAccessRoles(ctx echo.Context, repoId openapi_types.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}

func (h *APIHandler) CreateAccessRole(ctx echo.Context, repoId openapi_types.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Code:    intPtr(http.StatusNotImplemented),
		Message: strPtr("Not implemented"),
	})
}
