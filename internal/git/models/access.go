package git

import (
	"database/sql"
	"time"

	"github.com/gobwas/glob"
	"github.com/google/uuid"
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
	CreateUser(user *User) error
	EditUser(userid IDT, user *User) error
	DeleteUser(userid IDT) error
	GetUser(userid IDT) (User, error)
	GetUsers() ([]User, error)

	AddSshKey(userid IDT, key *SshKey) error
	DelSshKey(keyid IDT) error
	GetSshKeys(userid IDT) ([]SshKey, error)

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

func CreateUser(user *User) error {
	database, _ := sql.Open("sqlite3", "./access.db")
	statement, _ := database.Prepare("CREATE IF NOT EXISTS users ( id INTEGER PRIMARY KEY, name TEXT, email TEXT, created TIME, edited TIME . deleted TIME)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO users (id, name, email, created, edited, deleted) VALUES (?,?,?,?,?)")
	statement.Exec(user.ID, user.Name, user.Created, user.Edited, user.Deleted)
	return nil
}

func EditUser(userid IDT, user *User) error {
	database, _ := sql.Open("sqlite3", "./access.db")
	statement, _ := database.Prepare("UPDATE users SET name = ? , email = ? , created = ?, edited = ?, deleted = ? WHERE id = ?")
	statement.Exec(user.Name, user.Email, user.Created, user.Edited, user.Deleted, userid)
	return nil
}

func DeleteUser(userid IDT) error {
	database, _ := sql.Open("sqlite3", "./access.db")
	statement, _ := database.Prepare("DELETE FROM users WHERE id = ?")
	statement.Exec(userid)
	return nil
}

func GetUser(userid IDT) (User, error) {
	database, _ := sql.Open("sqlite3", "./access.db")
	row, _ := database.Query("SELECT * FORM users WHERE id = ?", userid)
	var user User
	for row.Next() {
		row.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
	}

	return user, nil
}

func GetUsers() ([]User, error) {
	database, _ := sql.Open("sqlite3", "./access.db")
	row, _ := database.Query("SELECT * FORM users ")
	var users []User
	for row.Next() {
		var user User
		row.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
		users = append(users, user)
	}

	return users, nil
}

func AddSshKey(userid IDT, key *SshKey) error {
	database, _ := sql.Open("sqlite3", "./access.db")
	statement, _ := database.Prepare("CREATE IF NOT EXISTS keys ( id INTEGER PRIMARY KEY, user_id INTEGER PRIMARY KEY, name TEXT, type INT  data VARCHAR(255), created TIME . deleted TIME)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO keys (id,user_id, name, type, data,  created, deleted) VALUES (?,?,?,?,?,?)")
	statement.Exec(key.ID, key.UserID, key.Name, key.Type, key.Data, key.Created, key.Deleted)
	return nil
}

func DelSshKey(keyid IDT) error {
	database, _ := sql.Open("sqlite3", "./access.db")
	statement, _ := database.Prepare("DELETE * FORM keys WHERE id = ?")
	statement.Exec(keyid)
	return nil
}

func GetSshKeys(userid IDT) ([]SshKey, error) {
	database, _ := sql.Open("sqlite3", "./access.db")
	row, _ := database.Query("SELECT * FROM keys WHERE user_id = ?", userid)
	var keys []SshKey
	for row.Next() {
		var key SshKey
		row.Scan(&key.ID, &key.UserID, &key.Name, &key.Type, &key.Data, &key.Created, &key.Deleted)
		keys = append(keys, key)
	}
	return keys, nil
}

func CreateRepo(repo *Repo) error {
	database, _ := sql.Open("sqlite3", "./access.db")
	statement, _ := database.Prepare("IF NOT EXISTS permitions (id INTEGER PRIMARY KEY, id_user INTEGER PRIMARY KEY, name TEXT, created TIME, deleted TIME")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO permitions (")
	return nil
}
