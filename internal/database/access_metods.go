package databases

import (
	"database/sql"
	"encoding/hex"
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
  RSA_SHA2-256 SSSSH_KEY_TYPE = iota + 2
  RSA_SHA2-512 SSSSH_KEY_TYPE = iota +3
  SSH_RSA  SSSSH_KEY_TYPE = iota +4
  ECDSA-SHA2-NISTP256 SSSSH_KEY_TYPE = iota +5
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
	UserID uuid.UUID
	RepoID uuid.UUID

	Branches *glob.Glob

	Created time.Time
	Deleted time.Time
}

type DB interface {
	CreateRepo(repo *Repo) error
	DeleteRepo(repoid IDT) error
	UpdateRepo(repoid IDT, repo *Repo) error
	GetRepo(repoid IDT) (*Repo, error)

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

var database *sql.DB

func init() {
	var err error
	database, err = sql.Open("sqlite3", "./migrations/access.db")
	if err != nil {
    logrus.Error("Error %s", err)
	}
	err = database.Ping()
	if err != nil {
    logrus.Error("Error %s", err)
	}
  logrus.Info("Database open")
}

func CreateUser(user *User) error {
	statement, err := database.Prepare("INSERT INTO users (id, name, email, created, edited, deleted) VALUES (?,?,?,?,?)")
	if err != nil {
    logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(hex.EncodeToString(user.ID), user.Name, user.Created, user.Edited, user.Deleted)
  logrus.Info("Create User", user)
	return nil
}

func EditUser(userid IDT, user *User) error {
	statement, err := database.Prepare("UPDATE users SET name = ? , email = ? , created = ?, edited = ?, deleted = ? WHERE id = ?")
	if err != nil {
    logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(user.Name, user.Email, user.Created, user.Edited, user.Deleted, hex.EncodeToString(userid))
  logrus.Info("Edit user", user)
	return nil
}

func DeleteUser(userid IDT) error {
	statement, err := database.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
    logrus.Error("Error %s", err)
		return err
	}
	 statement.Exec(userid)
  logrus.Info("Delete user", userid)
	return nil
}

func GetUser(userid IDT) (User, error) {
  var user User
	row, err := database.Query("SELECT * FORM users WHERE id = ?", userid)
	if err != nil {
    logrus.Error("Error %s", err)
		return user, err
	}
	var identif string
	row.Scan(&identif, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
	user.ID, err = hex.DecodeString(identif)
	if err != nil {
    logrus.Error("Error %s", err)
		return user, err
	}
  logrus.Trace("get info about user", user)
	return user, nil
}

func GetUsers() ([]User, error) {
  users := make([]User, 0, 16)
	row, err := database.Query("SELECT * FORM users ")
	if err != nil {
    logrus.Error("Error %s", err)
		return users, err
	}
	for row.Next() {
		var user User
		var identif string
		row.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
		user.ID, err = hex.DecodeString(identif)
		if err != nil {
    logrus.Error("Error %s", err)
			return users, err
		}
		users = append(users, user)
    logrus.Trace("get info about user", user)
	}

	return users, nil
}

func AddSshKey(userid IDT, key *SshKey) error {
	statement, err := database.Prepare("INSERT INTO keys (id,user_id, name, type, data,  created, deleted) VALUES (?,?,?,?,?,?)")
	if err != nil {
    logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(hex.EncodeToString(key.ID), hex.EncodeToString(key.UserID), key.Name, key.Type, key.Data, key.Created, key.Deleted)
  logrus.Info("Add SSH KEY", key, " ", userid)
	return nil
}

func DelSshKey(userid IDT, keyid IDT) error {
	statement, err := database.Prepare("DELETE * FORM keys WHERE id = ?, AND WHERE user_id = ?")
	if err != nil {
    logrus.Error("Error %s", err)
		return err
	}
	statement.Exec(keyid, userid)
  logrus.Info("Delete SSH key ", userid, " ", keyid)
	return nil
}

func GetSshKeys(userid IDT) ([]SshKey, error) {
  var keys []SshKey
	row, err := database.Query("SELECT * FROM keys WHERE user_id = ?", userid)
	if err != nil {
    logrus.Error("Error %s", err)
		return keys, err
	}
	for row.Next() {
		var key SshKey
		row.Scan(&key.ID, &key.UserID, &key.Name, &key.Type, &key.Data, &key.Created, &key.Deleted)
		keys = append(keys, key)
	}
  logrus.Info("Get SSH keys user ", userid )
	return keys, nil
}

func CreateRepo(repo *Repo) error {
	statement, err := database.Prepare("INSERT INTO permitions (")
	if err != nil {
    logrus.Error("Error %s", err)
		return err
	}
	return nil
}
