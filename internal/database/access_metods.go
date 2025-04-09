package database

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"time"

	"github.com/GoldenDeals/DepGit/internal/config"
	migrations "github.com/GoldenDeals/DepGit/internal/database/migrations"
	errors "github.com/GoldenDeals/DepGit/internal/share/error"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/gobwas/glob"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/sirupsen/logrus"
)

const (
	SSH_KEY_TYPE_RSA SSH_KEY_TYPE = iota + 1
	RSA_SHA2_256
	RSA_SHA2_512
	SSH_RSA
	ECDSA_SHA2_NISTP256
	// .... TODO: Добавить еще
)

// Logger for database access
var dbLogger = logger.New("db_access")

type (
	IDT          = uuid.UUID
	SSH_KEY_TYPE int
)

type User struct {
	ID uuid.UUID

	Name  string
	Email string

	Created time.Time
	Edited  time.Time
	Deleted time.Time
}

type SshKey struct {
	ID     uuid.UUID
	UserID uuid.UUID

	Name string
	Type SSH_KEY_TYPE
	Data []byte

	Created time.Time
	Deleted time.Time
}

type Repo struct {
	ID uuid.UUID

	Name string

	Created time.Time
	Edited  time.Time
	Deleted time.Time
}

type AccessRole struct {
	// PRIMARY KEY ( UserID, RepoID )
	RoleID uuid.UUID
	UserID uuid.UUID
	RepoID uuid.UUID

	Branches *glob.Glob

	Created time.Time
	Deleted time.Time
}

type Classes interface {
	New()
	Create()
}

type DB struct {
	db               *sql.DB
	migrationManager *migrations.Manager
	config           *config.Configuration
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Init(cfg *config.Configuration) error {
	var err error
	d.config = cfg
	d.db, err = sql.Open("sqlite3", cfg.GetDatabasePath())
	if err != nil {
		return err
	}
	err = d.db.Ping()
	if err != nil {
		return err
	}

	// Initialize the migration manager
	d.migrationManager = migrations.New(d.db)

	// Initialize the migrations table
	ctx := context.Background()
	err = d.migrationManager.Initialize(ctx)
	if err != nil {
		return err
	}

	// Check if custom migrations path is configured
	migrationsDir := d.config.DB.MigrationsPath
	if migrationsDir == "" {
		migrationsDir = filepath.Join(filepath.Dir(cfg.GetDatabasePath()), "migrations")
	}

	err = d.migrationManager.LoadFromDir(ctx, migrationsDir)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("error", err).
			Warn("Failed to load migrations")
	}

	// Apply migrations
	err = d.migrationManager.Apply(ctx)
	if err != nil {
		return err
	}

	dbLogger.
		WithContext(ctx).
		WithField("database", cfg.GetDatabasePath()).
		Info("Database open")
	return nil
}

func NewUser(name, email string) User {
	return User{
		ID:      uuid.New(),
		Name:    name,
		Email:   email,
		Created: time.Now(),
	}
}

func (d *DB) CreateUser(ctx context.Context, user *User) error {
	if user == nil {
		return errors.ERR_BAD_DATA
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	// TODO: More validations
	if user.ID == uuid.Nil || user.Name == "" || user.Email == "" || !strings.ContainsRune(user.Email, '@') {
		return errors.ERR_BAD_DATA
	}

	// TODO: Error check or recovery
	userid := user.ID.String()
	var n int
	row := d.db.QueryRow("SELECT COUNT(id) FROM users WHERE id = ?", userid)
	row.Scan(&n)
	if n > 0 {
		return errors.ERR_ALREADY_EXISTS
	}
	statement, err := d.db.Prepare("INSERT INTO users (id, name, email, created, edited, deleted) VALUES (?,?,?,?,?,?)")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error creating user")

		return err
	}
	_, err = statement.Exec(userid, user.Name, user.Email, time.Now().Format(time.DateTime), user.Edited.Format(time.DateTime), user.Deleted.Format(time.DateTime))
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error creating user")

		return err
	}
	logrus.Trace("Create User", userid, user.Name)
	return nil
}

func (d *DB) EditUser(ctx context.Context, userid IDT, user *User) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if user == nil || user.ID == uuid.Nil || user.Name == "" || user.Email == "" || !strings.ContainsRune(user.Email, '@') {
		return errors.ERR_BAD_DATA
	}

	statement, err := d.db.Prepare("UPDATE users SET name = ? , email = ? , created = ?, edited = ?, deleted = ? WHERE id = ?")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error edit user")
		return err
	}
	_, err = statement.Exec(user.Name, user.Email, user.Created.Format(time.DateTime), time.Now().Format(time.DateTime), user.Deleted.Format(time.DateTime), userid.String())
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error edit user")
		return err
	}
	logrus.Trace("Edit user", userid, user.Name)
	return nil
}

func (d *DB) DeleteUser(ctx context.Context, userid IDT) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if userid == uuid.Nil {
		return errors.ERR_BAD_DATA
	}

	var n int
	row := d.db.QueryRow("SELECT COUNT(id) FROM users WHERE id = ?", userid)
	err := row.Scan(&n)
	if err != nil {
		dbLogger.WithContext(ctx).WithError(err).Warn("error checking user existence")
		return err
	}

	if n == 0 {
		return errors.ERR_NOT_FOUND
	}

	statement, err := d.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error delete user")
		return err
	}
	_, err = statement.Exec(userid)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error delete user")
		return err
	}
	logrus.Trace("Delete user", userid)
	return nil
}

func (d *DB) GetUser(ctx context.Context, userid IDT) (User, error) {
	var user User
	if err := ctx.Err(); err != nil {
		return user, err
	}
	if userid == uuid.Nil {
		return user, errors.ERR_BAD_DATA
	}
	var n int
	row := d.db.QueryRow("SELECT COUNT(id) FROM users WHERE id = ?", userid)
	err := row.Scan(&n)
	if err != nil {
		dbLogger.WithContext(ctx).WithError(err).Warn("error checking user existence")
		return user, err
	}

	if n == 0 {
		return user, errors.ERR_NOT_FOUND
	}

	row = d.db.QueryRow("SELECT id, name, email, created, edited, deleted FROM users WHERE id = ?", userid)
	var idStr string
	err = row.Scan(&idStr, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error get user")
		return user, err
	}

	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error parsing user id")
		return user, err
	}

	logrus.Trace("get info about user", user)
	return user, nil
}

func (d *DB) GetUsers(ctx context.Context) ([]User, error) {
	users := make([]User, 0, 16)
	if err := ctx.Err(); err != nil {
		return users, err
	}

	rows, err := d.db.Query("SELECT id, name, email, created, edited, deleted FROM users")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithError(err).
			Warn("error fetching users")
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		var idStr string
		err = rows.Scan(&idStr, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
		if err != nil {
			dbLogger.
				WithContext(ctx).
				WithError(err).
				Warn("error scanning user row")
			return users, err
		}

		user.ID, err = uuid.Parse(idStr)
		if err != nil {
			dbLogger.
				WithContext(ctx).
				WithField("user_id", idStr).
				WithError(err).
				Warn("error parsing user id")
			return users, err
		}

		users = append(users, user)
		logrus.Trace("get info about user", user.ID, user.Name)
	}

	if err = rows.Err(); err != nil {
		dbLogger.
			WithContext(ctx).
			WithError(err).
			Warn("error iterating user rows")
		return users, err
	}

	return users, nil
}

func NewSShKey(name string, key_type SSH_KEY_TYPE, data []byte) SshKey {
	return SshKey{
		ID:      uuid.New(),
		Name:    name,
		Type:    key_type,
		Data:    data,
		Created: time.Now(),
	}
}

func (d *DB) AddSshKey(ctx context.Context, userid IDT, key *SshKey) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// Set the UserID to the provided userid
	key.UserID = userid

	if userid == uuid.Nil || key.ID == uuid.Nil || key.Name == "" || key.Type > 5 {
		return errors.ERR_BAD_DATA
	}
	// Check if key already exists
	row := d.db.QueryRow("SELECT COUNT(id) FROM keys WHERE user_id = ? AND data = ?", userid, key.Data)
	var n int
	var err error
	err = row.Scan(&n)
	if err != nil {
		dbLogger.WithContext(ctx).WithError(err).Warn("error checking key existence")
		return err
	}

	if n > 0 {
		return errors.ERR_ALREADY_EXISTS
	}
	statement, err := d.db.Prepare("INSERT INTO keys (id, user_id, name, type, data, created, deleted) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithField("key", key.ID).
			WithError(err).
			Warn("error add ssh key")
		return err
	}
	_, err = statement.Exec(key.ID.String(), key.UserID.String(), key.Name, key.Type, key.Data, time.Now().Format(time.DateTime), key.Deleted.Format(time.DateTime))
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithField("key", key.ID).
			WithError(err).
			Warn("error add ssh key")
		return err
	}
	logrus.Trace("Add SSH KEY", key, userid)
	return nil
}

func (d *DB) DeleteSshKey(ctx context.Context, keyid IDT) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	row := d.db.QueryRow("SELECT COUNT(id) FROM  keys  WHERE id = ?", keyid)
	var n int
	row.Scan(&n)
	if n == 0 {
		return errors.ERR_NOT_FOUND
	}
	statement, err := d.db.Prepare("DELETE FROM keys WHERE id = ?")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("key_id", keyid).
			WithError(err).
			Warn("error add ssh key")
		return err
	}
	_, err = statement.Exec(keyid)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("key_id", keyid).
			WithError(err).
			Warn("error add ssh key")
		return err
	}
	logrus.Trace("Delete SSH key ", keyid)
	return nil
}

func (d *DB) GetSshKeys(ctx context.Context, userid IDT) ([]SshKey, error) {
	keys := make([]SshKey, 0, 16)
	if err := ctx.Err(); err != nil {
		return keys, err
	}
	if userid == uuid.Nil {
		return keys, errors.ERR_BAD_DATA
	}
	len := d.db.QueryRow("SELECT COUNT(id) FROM  keys  WHERE user_id = ?", userid)
	var n int
	len.Scan(&n)
	if n == 0 {
		return keys, errors.ERR_NOT_FOUND
	}
	row, err := d.db.Query("SELECT * FROM keys WHERE user_id = ?", userid)
	if err != nil { //
		dbLogger.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error get ssh key")
		return keys, err
	}
	for row.Next() {
		var key SshKey
		err = row.Scan(&key.ID, &key.UserID, &key.Name, &key.Type, &key.Data, &key.Created, &key.Deleted)
		if err != nil {
			dbLogger.
				WithContext(ctx).
				WithField("user_id", userid).
				WithError(err).
				Warn("error get ssh key")
			return keys, err
		}
		keys = append(keys, key)
	}
	logrus.Trace("Get SSH keys user ", userid)
	return keys, nil
}

func NewRepo(name string) Repo {
	return Repo{
		ID:   uuid.New(),
		Name: name,
	}
}

func (d *DB) CreateRepo(ctx context.Context, repo *Repo) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if repo.ID == uuid.Nil || repo.Name == "" {
		return errors.ERR_BAD_DATA
	}
	statement, err := d.db.Prepare("INSERT INTO permitions (id, name, created, edited, deleted) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("repoid", repo.ID).
			WithError(err).
			Warn("error create repo")
		return err
	}
	defer statement.Close()

	// Use the repo ID directly
	_, err = statement.Exec(
		repo.ID.String(),
		repo.Name,
		time.Now().Format(time.DateTime),
		repo.Edited.Format(time.DateTime),
		repo.Deleted.Format(time.DateTime),
	)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("repoid", repo.ID).
			WithError(err).
			Warn("error create repo")
		return err
	}
	logrus.Trace("Create Repository", repo.ID)
	return nil
}

func (d *DB) DeleteRepo(ctx context.Context, repoid IDT) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if repoid == uuid.Nil {
		return errors.ERR_BAD_DATA
	}
	statement, err := d.db.Prepare("DELETE FROM permitions WHERE id = ?")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error delete repo")
		return err
	}
	_, err = statement.Exec(repoid)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error delete repo")
		return err
	}
	logrus.Trace("Deleted.Format(time.DateTime) Repositor", repoid)
	return nil
}

func (d *DB) UpdateRepo(ctx context.Context, repoid IDT, repo Repo) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if repoid == uuid.Nil || repo.ID == uuid.Nil || repo.Name == "" {
		return errors.ERR_BAD_DATA
	}

	statement, err := d.db.Prepare("UPDATE permitions SET name = ? , created = ?, deleted = ? WHERE id = ?")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error update repo")
		return err
	}
	_, err = statement.Exec(repo.Name, repo.Created.Format(time.DateTime), repo.Deleted.Format(time.DateTime), repoid.String())
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error update repo")
		return err
	}
	logrus.Trace("Edit Repositor ", repo.ID)
	return nil
}

func (d *DB) GetRepo(ctx context.Context, repoid IDT) (Repo, error) {
	var repo Repo
	if err := ctx.Err(); err != nil {
		return repo, err
	}
	if repoid == uuid.Nil {
		return repo, errors.ERR_BAD_DATA
	}

	var n int
	row := d.db.QueryRow("SELECT COUNT(id) FROM permitions WHERE id = ?", repoid)
	err := row.Scan(&n)
	if err != nil {
		dbLogger.WithContext(ctx).WithError(err).Warn("error checking repo existence")
		return repo, err
	}

	if n == 0 {
		return repo, errors.ERR_NOT_FOUND
	}

	row = d.db.QueryRow("SELECT id, name, created, edited, deleted FROM permitions WHERE id = ?", repoid)
	var idStr string
	err = row.Scan(&idStr, &repo.Name, &repo.Created, &repo.Edited, &repo.Deleted)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error get repo")
		return repo, err
	}

	repo.ID, err = uuid.Parse(idStr)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("repoid", idStr).
			WithError(err).
			Warn("error parsing repo id")
		return repo, err
	}

	logrus.Trace("get info about Repositor ", repo.ID)
	return repo, nil
}

func (d *DB) GetRepos(ctx context.Context) ([]Repo, error) {
	repos := make([]Repo, 0, 16)
	if err := ctx.Err(); err != nil {
		return repos, err
	}

	rows, err := d.db.Query("SELECT id, name, created, edited, deleted FROM permitions")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithError(err).
			Warn("error fetching repos")
		return repos, err
	}
	defer rows.Close()

	for rows.Next() {
		var repo Repo
		var idStr string
		err = rows.Scan(&idStr, &repo.Name, &repo.Created, &repo.Edited, &repo.Deleted)
		if err != nil {
			dbLogger.
				WithContext(ctx).
				WithError(err).
				Warn("error scanning repo row")
			return repos, err
		}

		repo.ID, err = uuid.Parse(idStr)
		if err != nil {
			dbLogger.
				WithContext(ctx).
				WithField("repo_id", idStr).
				WithError(err).
				Warn("error parsing repo id")
			return repos, err
		}

		repos = append(repos, repo)
		logrus.Trace("get info about Repositor", repo.ID)
	}

	if err = rows.Err(); err != nil {
		dbLogger.
			WithContext(ctx).
			WithError(err).
			Warn("error iterating repo rows")
		return repos, err
	}

	return repos, nil
}

func (d *DB) CreateAccessRole(ctx context.Context, ar *AccessRole) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if ar.UserID == uuid.Nil || ar.RepoID == uuid.Nil || ar.RoleID == uuid.Nil {
		return errors.ERR_BAD_DATA
	}
	statement, err := d.db.Prepare("INSERT INTO roles (role_id, user_id, rep_id, branch, created, deleted) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("role", ar.RoleID).
			WithError(err).
			Warn("error create role")
		return err
	}
	defer statement.Close()

	// Use the role's ID directly
	_, err = statement.Exec(
		ar.RoleID.String(),
		ar.UserID.String(),
		ar.RepoID.String(),
		ar.Branches,
		ar.Created.Format(time.DateTime),
		ar.Deleted.Format(time.DateTime),
	)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("role", ar.RoleID).
			WithError(err).
			Warn("error create role")
		return err
	}
	logrus.Trace("Create Role", ar.RoleID)
	return nil
}

func (d *DB) EditAccessRole(ctx context.Context, roleid IDT, ar *AccessRole) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if roleid == uuid.Nil || ar.UserID == uuid.Nil || ar.RepoID == uuid.Nil || ar.RoleID == uuid.Nil {
		return errors.ERR_BAD_DATA
	}
	statement, err := d.db.Prepare("UPDATE roles SET user_id = ?, rep_id = ?, branch = ?, created = ?, deleted = ? WHERE role_id = ?")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error edit role")
		return err
	}
	_, err = statement.Exec(ar.UserID.String(), ar.RepoID.String(), ar.Branches, ar.Created.Format(time.DateTime), ar.Deleted.Format(time.DateTime), ar.RoleID.String())
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("role ", ar.RoleID).
			WithError(err).
			Warn("error edit role")
		return err
	}
	logrus.Trace("Edit role", ar.RoleID)
	return nil
}

func (d *DB) DeleteAccessRole(ctx context.Context, roleid IDT) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if roleid == uuid.Nil {
		return errors.ERR_BAD_DATA
	}
	statement, err := d.db.Prepare("DELETE FROM roles WHERE role_id = ?")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error deleted role")
		return err
	}
	_, err = statement.Exec(roleid)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error deleted role")
		return err
	}
	logrus.Trace("Deleted.Format(time.DateTime) Role", roleid)
	return nil
}

func (d *DB) GetAccessRole(ctx context.Context, roleid IDT) (AccessRole, error) {
	var role AccessRole
	if err := ctx.Err(); err != nil {
		return role, err
	}
	if roleid == uuid.Nil {
		return role, errors.ERR_BAD_DATA
	}

	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM roles WHERE role_id = ?", roleid.String()).Scan(&count)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("role_id", roleid).
			WithError(err).
			Warn("error checking role existence")
		return role, err
	}

	if count == 0 {
		return role, errors.ERR_NOT_FOUND
	}

	// Variables to store the data
	var roleIDStr, userIDStr, repoIDStr string

	// Query the role
	row := d.db.QueryRow("SELECT role_id, user_id, rep_id, branch, created, deleted FROM roles WHERE role_id = ?", roleid.String())

	// Scan the row into variables
	err = row.Scan(&roleIDStr, &userIDStr, &repoIDStr, &role.Branches, &role.Created, &role.Deleted)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("role_id", roleid).
			WithError(err).
			Warn("error scanning role")
		return role, err
	}

	// Parse the role ID
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("role_id_str", roleIDStr).
			WithError(err).
			Warn("error parsing role ID")
		return role, err
	}
	role.RoleID = roleID

	// Parse the user ID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("user_id_str", userIDStr).
			WithError(err).
			Warn("error parsing user ID")
		return role, err
	}
	role.UserID = userID

	// Parse the repo ID
	repoID, err := uuid.Parse(repoIDStr)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithField("repo_id_str", repoIDStr).
			WithError(err).
			Warn("error parsing repo ID")
		return role, err
	}
	role.RepoID = repoID

	dbLogger.WithContext(ctx).Debug("Retrieved role information")
	return role, nil
}

func (d *DB) GetAccessRoles(ctx context.Context) ([]AccessRole, error) {
	var roles []AccessRole
	if err := ctx.Err(); err != nil {
		return roles, err
	}

	rows, err := d.db.Query("SELECT role_id, user_id, rep_id, branch, created, deleted FROM roles")
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithError(err).
			Warn("error get roles")
		return roles, err
	}
	defer rows.Close()

	for rows.Next() {
		var role AccessRole
		var roleIDStr, userIDStr, repoIDStr string

		err = rows.Scan(&roleIDStr, &userIDStr, &repoIDStr, &role.Branches, &role.Created, &role.Deleted)
		if err != nil {
			dbLogger.
				WithContext(ctx).
				WithError(err).
				Warn("error scanning role")
			return roles, err
		}

		// Parse IDs
		role.RoleID, err = uuid.Parse(roleIDStr)
		if err != nil {
			dbLogger.
				WithContext(ctx).
				WithField("role_id_str", roleIDStr).
				WithError(err).
				Warn("error parsing role ID")
			return roles, err
		}

		role.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			dbLogger.
				WithContext(ctx).
				WithField("user_id_str", userIDStr).
				WithError(err).
				Warn("error parsing user ID")
			return roles, err
		}

		role.RepoID, err = uuid.Parse(repoIDStr)
		if err != nil {
			dbLogger.
				WithContext(ctx).
				WithField("repo_id_str", repoIDStr).
				WithError(err).
				Warn("error parsing repo ID")
			return roles, err
		}

		roles = append(roles, role)
		logrus.Trace("get info about Role", role.RoleID)
	}

	if err = rows.Err(); err != nil {
		dbLogger.
			WithContext(ctx).
			WithError(err).
			Warn("error iterating role rows")
		return roles, err
	}

	return roles, nil
}

func (d *DB) UserByKey(ctx context.Context, key []byte) (User, error) {
	var user User
	var err error
	if err := ctx.Err(); err != nil {
		return user, err
	}
	if key == nil {
		return user, errors.ERR_BAD_DATA
	}
	row := d.db.QueryRow("SELECT userid.keys, users.name, users.email , users.created, users.edited, users.deleted FROM keys INNER JOIN users ON users.userid = keys.userid", key) // что делает и что выводит ?
	var answer string
	err = row.Scan(&answer, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithError(err).
			WithField("key", key).
			Warn("error get user by key")
		return user, err
	}
	logrus.Trace("get info user keys ", user, key)
	return user, nil
}

func (d *DB) CheckPermissions(ctx context.Context, userid IDT, repoid IDT, branch string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if userid == uuid.Nil || repoid == uuid.Nil {
		return false, errors.ERR_BAD_DATA
	}

	// No branch matching - just check if the user has any role for the repo
	query := "SELECT COUNT(*) FROM roles WHERE user_id = ? AND rep_id = ?"
	args := []interface{}{userid.String(), repoid.String()}

	// Add branch condition if provided
	if branch != "" {
		query += " AND branch = ?"
		args = append(args, branch)
	}

	var count int
	err := d.db.QueryRow(query, args...).Scan(&count)

	if err != nil {
		dbLogger.
			WithContext(ctx).
			WithError(err).
			WithField("user_id", userid).
			WithField("repo_id", repoid).
			Warn("error check permissions")
		return false, err
	}

	// If we found any roles, the user has permission
	return count > 0, nil
}
