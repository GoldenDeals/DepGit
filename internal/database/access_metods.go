package database

import (
	"database/sql"
	"time"

	"github.com/gobwas/glob"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IDT = uuid.UUID

type User struct {
	ID uuid.UUID

	Name  string
	Email string

	Created time.Time
	Edited  time.Time
	Deleted time.Time
}

type SSH_KEY_TYPE int

const (
	SSH_KEY_TYPE_RSA SSH_KEY_TYPE = iota + 1
	// RSA_SHA2_256 SSSSH_KEY_TYPE
	// RSA_SHA2_512 SSSSH_KEY_TYPE
	// SSH_RSA  SSSSH_KEY_TYPE
	// ECDSA_SHA2NISTP256 SSSSH_KEY_TYPE
	// .... TODO: Добавить еще
)

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

type DBa interface {
	CreateRepo(repo *Repo) error
	DeleteRepo(repoid IDT) error
	UpdateRepo(repoid IDT, repo *Repo) error
	GetRepo(repoid IDT) (*Repo, error)
	GetRepos() ([]Repo, error)

	// Это функции добавляющие пользователя в репозиторий
	CreateAccessRole(ar *AccessRole) error
	EditAccessRole(roleid IDT, ar *AccessRole)
	DeleteAccessRole(roleid IDT) error
	GetAccessRoles(repoid IDT) ([]AccessRole, error)
}

type AccessManager interface {
	UserByKey(key []byte) (*User, error)
	CheckPermissions(userid IDT, repoid IDT, branch string) (bool, error)
}

type DB struct {
	db *sql.DB
}

func (d *DB) init() (err error) {
	d.db, err = sql.Open("sqlite3", "./migrations/access.db")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	err = d.db.Ping()
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	logrus.Trace("Database access.db open")
	return nil
}

func (d *DB) CreateUser(user *User) (err error) {
	d.init()
	statement, err := d.db.Prepare("INSERT INTO users (id, name, email, created, edited, deleted) VALUES (?,?,?,?,?)")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(user.ID.String(), user.Name, user.Created, user.Edited, user.Deleted)
	logrus.Trace("Create User", user)
	return nil
}

func (d *DB) EditUser(userid IDT, user *User) error {
	statement, err := d.db.Prepare("UPDATE users SET name = ? , email = ? , created = ?, edited = ?, deleted = ? WHERE id = ?")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(user.Name, user.Email, user.Created, user.Edited, user.Deleted, userid.String())
	logrus.Trace("Edit user", user)
	return nil
}

func (d *DB) DeleteUser(userid IDT) error {
	statement, err := d.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(userid)
	logrus.Trace("Delete user", userid)
	return nil
}

func (d *DB) GetUser(userid IDT) (User, error) {
	var user User
	row, err := d.db.Query("SELECT * FORM users WHERE id = ?", userid)
	if err != nil {
		logrus.Error("Error %s", err)
		return user, err
	}
	var identif string
	row.Scan(&identif, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
	user.ID, err = uuid.FromString(identif)
	if err != nil {
		logrus.Error("Error %s", err)
		return user, err
	}
	logrus.Trace("get info about user", user)
	return user, nil
}

func (d *DB) GetUsers() ([]User, error) {
	users := make([]User, 0, 16)
	row, err := d.db.Query("SELECT * FORM users ")
	if err != nil {
		logrus.Error("Error %s", err)
		return users, err
	}
	for row.Next() {
		var user User
		var identif string
		row.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
		user.ID, err = uuid.FromString(identif)
		if err != nil {
			logrus.Error("Error %s", err)
			return users, err
		}
		users = append(users, user)
		logrus.Trace("get info about user", user)
	}

	return users, nil
}

func (d *DB) AddSshKey(userid IDT, key *SshKey) error {
	statement, err := d.db.Prepare("INSERT INTO keys (id,user_id, name, type, data,  created, deleted) VALUES (?,?,?,?,?,?)")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(key.ID.String(), key.UserID.String(), key.Name, key.Type, key.Data, key.Created, key.Deleted)
	logrus.Trace("Add SSH KEY", key, " ", userid)
	return nil
}

func (d *DB) DelSshKey(userid IDT, keyid IDT) error {
	statement, err := d.db.Prepare("DELETE * FORM keys WHERE id = ?, AND WHERE user_id = ?")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(keyid, userid)
	logrus.Trace("Delete SSH key ", userid, " ", keyid)
	return nil
}

func (d *DB) GetSshKeys(userid IDT) ([]SshKey, error) {
	keys := make([]SshKey, 0, 16)
	row, err := d.db.Query("SELECT * FROM keys WHERE user_id = ?", userid)
	if err != nil {
		logrus.Error("Error %s", err)
		return keys, err
	}
	for row.Next() {
		var key SshKey
		row.Scan(&key.ID, &key.UserID, &key.Name, &key.Type, &key.Data, &key.Created, &key.Deleted)
		keys = append(keys, key)
	}
	logrus.Trace("Get SSH keys user ", userid)
	return keys, nil
}

func (d *DB) CreateRepo(repo *Repo) error {
	statement, err := d.db.Prepare("INSERT INTO permitions (id, name, created, deleted) VALUES (?, ?, ?, ?)")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(repo.ID, repo.Name, repo.Created, repo.Deleted)
	logrus.Trace("Create Repositor", repo)
	return nil
}

func (d *DB) DeleteRepo(repoid IDT) error {
	statement, err := d.db.Prepare("DELETE * FORM permitions WHERE id = ?")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(repoid)
	logrus.Trace("Deleted Repositor", repoid)
	return nil
}

func (d *DB) UpdateRepo(repoid IDT, repo Repo) error {
	statement, err := d.db.Prepare("UPDATE permitions SET name = ? , created = ?, deleted = ? WHERE id = ?")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(repo.Name, repo.Created, repo.Deleted, repo.ID.String())
	logrus.Trace("Edit Repositor ", repo)
	return nil
}

func (d *DB) GetRepo(repoid IDT) (*Repo, error) {
	var repo Repo
	row, err := d.db.Query("SELECT * FORM permitions WHERE id = ?", repoid)
	if err != nil {
		logrus.Error("Error %s", err)
		return &repo, err
	}
	var identif string
	row.Scan(&identif, &repo.Name, &repo.Created, &repo.Deleted)
	repo.ID, err = uuid.FromString(identif)
	if err != nil {
		logrus.Error("Error %s", err)
		return &repo, err
	}
	logrus.Trace("get info about Repositor ", repo)
	return &repo, nil
}

func (d *DB) GetRepos() ([]Repo, error) {
	var err error
	repos := make([]Repo, 0, 16)
	row, err := d.db.Query("SELECT * FORM permitions ")
	if err != nil {
		logrus.Error("Error %s", err)
		return repos, err
	}
	for row.Next() {
		var repo Repo
		var identif string
		row.Scan(&identif, &repo.Name, &repo.Created, &repo.Deleted)
		repo.ID, err = uuid.FromString(identif)
		if err != nil {
			logrus.Error("Error %s", err)
			return repos, err
		}
		repos = append(repos, repo)
		logrus.Trace("get info about Repositor", repo)
	}

	return repos, nil
}

func (d *DB) CreateAccessRole(ar *AccessRole) error {
	statement, err := d.db.Prepare("INSERT INTO roles (role_id ,user_id, rep_id , branch, created, deleted) VALUES (?,?, ?,?, ?, ?)")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(ar.RoleID, ar.UserID, ar.RepoID, ar.Branches, ar.Created, ar.Deleted)
	logrus.Trace("Create Repositor", ar)
	return nil
}

func (d *DB) EditAccessRole(roleid IDT, ar *AccessRole) error {
	statement, err := d.db.Prepare("UPDATE roles SET user_id = ?, rep_id = ?, branch = ?, created = ?, deleted = ? WHERE role_id = ?")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(ar.UserID.String(), ar.RepoID.String(), ar.Branches, ar.Created, ar.Deleted, ar.RoleID.String())
	logrus.Trace("Edit role", ar)
	return nil
}

func (d *DB) DeleteAccessRole(roleid IDT) error {
	statement, err := d.db.Prepare("DELETE * FORM roles WHERE role_id = ?")
	if err != nil {
		logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(roleid)
	logrus.Trace("Deleted Role", roleid)
	return nil
}

func (d *DB) GetAccessRole(roleid IDT) (AccessRole, error) {
	var role AccessRole
	row, err := d.db.Query("SELECT * FORM roles WHERE id = ?", roleid)
	if err != nil {
		logrus.Error("Error %s", err)
		return role, err
	}
	var identifRole, identifUser, identifReposit string
	row.Scan(&identifRole, &identifUser, &identifReposit, &role.Branches, &role.Created, &role.Deleted)
	role.RoleID, err = uuid.FromString(identifRole)
	role.UserID, err = uuid.FromString(identifUser)
	role.RepoID, err = uuid.FromString(identifReposit)
	if err != nil {
		logrus.Error("Error %s", err)
		return role, err
	}
	logrus.Trace("get info about Role ", role)
	return role, nil
}

func (d *DB) GetAccessRoles() ([]AccessRole, error) {
	var roles []AccessRole
	row, err := d.db.Query("SELECT * FORM roles ")
	if err != nil {
		logrus.Error("Error %s", err)
		return roles, err
	}
	for row.Next() {
		var identifRole, identifUser, identifReposit string
		var role AccessRole
		row.Scan(&identifRole, &identifUser, &identifReposit, &role.Branches, &role.Created, &role.Deleted)
		role.RoleID, err = uuid.FromString(identifRole)
		role.UserID, err = uuid.FromString(identifUser)
		role.RepoID, err = uuid.FromString(identifReposit)
		if err != nil {
			logrus.Error("Error %s", err)
			return roles, err
		}
		logrus.Trace("get info about Role ", role)
		roles = append(roles, role)
	}
	return roles, nil
}

func (d *DB) UserByKey(key []byte) (*User, error) {
	err := d.init()
	if err != nil {
		logrus.Error("Error %s", err)
		return nil, err
	}
	// что делает и что выводит ?
	return nil, nil
}

func (d *DB) CheckPermissions(userid IDT, repoid IDT, branch string) (bool, error) {
	err := d.init()
	if err != nil {
		logrus.Error("Error %s", err)
		return false, err
	}
	row, err := d.db.Query("SELECT roleid FROM roles WHERE userid = ? AND WHERE repoid = ? AND WHERE branch = ?", userid.String(), repoid.String(), branch)
	if err != nil {
		logrus.Error("Error %s", err)
		return false, err
	}
	var roleid string
	err = row.Scan(&roleid)
	if err != nil {
		logrus.Error("Error %s", err)
		return false, err
	}
	return true, err
}
