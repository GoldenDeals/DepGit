package git

import (
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
	Edited  time.Time
	Deleted time.Time
}

type AccessRole struct {
	// PRIMARY KEY ( UserID, RepoID )
	UserID uuid.UUID
	RepoID uuid.UUID

	Branches *glob.Glob

	Created time.Time
	Edited  time.Time
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
