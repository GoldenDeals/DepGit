package database

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/GoldenDeals/DepGit/internal/share/errors"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/gobwas/glob"
	"github.com/google/uuid"
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

var log = logger.New("db")

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
	db *sql.DB
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Init(dbname string) error {
	var err error
	d.db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	err = d.db.Ping()
	if err != nil {
		return err
	}
	log.Debugf("Database %s open", dbname)
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
	if err := ctx.Err(); err != nil {
		return err
	}

	// TODO: More validations
	if user == nil || user.ID == uuid.Nil || user.Name == "" || user.Email == "" || !strings.ContainsRune(user.Email, '@') {
		return errors.ERR_BAD_DATA
	}

	// TODO: Error check or recovery
	userid := user.ID.String()
	var n int
	row := d.db.QueryRow("SELECT COUNT(user_id) FROM  users WHERE userid = ?", userid)
	row.Scan(&n)
	if n > 0 {
		return errors.ERR_ALREADY_EXISTS
	}
	statement, err := d.db.Prepare("INSERT INTO users (id, name, email, created, edited, deleted) VALUES (?,?,?,?,?)")
	if err != nil {
		log.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error creating user")

		return err
	}
	_, err = statement.Exec(userid, user.Name, time.Now().Format(time.DateTime), user.Edited.Format(time.DateTime), user.Deleted.Format(time.DateTime))
	if err != nil {
		log.
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
		log.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error edit user")
		return err
	}
	_, err = statement.Exec(user.Name, user.Email, user.Created.Format(time.DateTime), time.Now().Format(time.DateTime), user.Deleted.Format(time.DateTime), userid.String())
	if err != nil {
		log.
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

	if len(userid) != 36 {
		return errors.ERR_BAD_DATA
	}

	var n int
	row := d.db.QueryRow("SELECT COUNT(userid) FROM users WHERE userid = ?", userid)
	row.Scan(&n)
	if n == 0 {
		return errors.ERR_NOT_FOUND
	}

	statement, err := d.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		log.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error delete user")
		return err
	}
	_, err = statement.Exec(userid)
	if err != nil {
		log.
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
	var err error
	var identif string
	var id IDT
	if err := ctx.Err(); err != nil {
		return user, err
	}
	if len(userid) != 36 {
		return errors.ERR_BAD_DATA
	}
	var n int
	row := d.db.QueryRow("SELECT COUNT(userid) FROM users WHERE userid = ?", userid)
	row.Scan(&n)
	if n == 0 {
		return errors.ERR_NOT_FOUND
	}
	row = d.db.QueryRow("SELECT * FORM users WHERE id = ?", userid)
	err = row.Scan(&identif, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
	if err != nil {
		log.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error get user")
		return user, err
	}
	err = id.UnmarshalText([]byte(identif))
	user.ID = id
	if err != nil {
		log.
			WithContext(ctx).
			WithField("user_id", userid).
			WithError(err).
			Warn("error get user")
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
	row, err := d.db.Query("SELECT * FORM users ")
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			Warn("error gets user")
		return users, err
	}
	for row.Next() {
		var user User
		var identif string
		var id IDT
		row.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
		err = id.UnmarshalText([]byte(identif))
		user.ID = id
		if err != nil {
			log.
				WithContext(ctx).
				WithField("user_id", identif).
				WithError(err).
				Warn("error get user")
			return users, err
		}
		users = append(users, user)
		logrus.Trace("get info about user", user.ID, user.Name)
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
	if userid == uuid.Nil || key.ID == uuid.Nil || key.Name == "" || key.Type > 5 || key.UserID == uuid.Nil {
		return errors.ERR_BAD_DATA
	}
	row := d.db.QueryRow("SELECT COUNT(id) FROM  keys WHERE userid = ? AND WHERE key = ?", userid, key.Data)
	var n int
	row.Scan(&n)
	if n > 0 {
		return errors.ERR_ALREADY_EXISTS
	}
	var err error
	statement, err := d.db.Prepare("INSERT INTO keys (id,user_id, name, type, key,  created, deleted) VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.
			WithContext(ctx).
			WithField("user_id", userid).
			WithField("key", key.ID).
			WithError(err).
			Warn("error add ssh key")
		return err
	}
	_, err = statement.Exec(key.ID.String(), key.UserID.String(), key.Name, key.Type, key.Data, time.Now().Format(time.DateTime), key.Deleted.Format(time.DateTime))
	if err != nil {
		log.
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
	statement, err := d.db.Prepare("DELETE * FORM keys WHERE id = ?")
	if err != nil {
		log.
			WithContext(ctx).
			WithField("key_id", keyid).
			WithError(err).
			Warn("error add ssh key")
		return err
	}
	_, err = statement.Exec(keyid)
	if err != nil {
		log.
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
		return errors.ERR_BAD_DATA
	}
	len := d.db.QueryRow("SELECT COUNT(id) FROM  keys  WHERE user_id = ?", userid)
	var n int
	len.Scan(&n)
	if n == 0 {
		return errors.ERR_NOT_FOUND
	}
	row, err := d.db.Query("SELECT * FROM keys WHERE user_id = ?", userid)
	if err != nil { //
		log.
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
			log.
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
	statement, err := d.db.Prepare("INSERT INTO permitions (id, name, created, edited ,deleted) VALUES (?, ?,?, ?, ?)")
	if err != nil {
		log.
			WithContext(ctx).
			WithField("repoid", repo.ID).
			WithError(err).
			Warn("error create repo")
		return err
	}
	repoid := uuid.New().String()
	_, err = statement.Exec(repoid, repo.Name, time.Now().Format(time.DateTime), repo.Edited, repo.Deleted.Format(time.DateTime))
	if err != nil {
		log.
			WithContext(ctx).
			WithField("repoid", repo.ID).
			WithError(err).
			Warn("error create repo")
		return err
	}
	logrus.Trace("Create Repositor", repo.ID)
	return nil
}

func (d *DB) DeleteRepo(ctx context.Context, repoid IDT) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if repoid == uuid.Nil {
		return errors.ERR_BAD_DATA
	}
	statement, err := d.db.Prepare("DELETE * FORM permitions WHERE id = ?")
	if err != nil {
		log.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error delete repo")
		return err
	}
	_, err = statement.Exec(repoid)
	if err != nil {
		log.
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
		log.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error update repo")
		return err
	}
	_, err = statement.Exec(repo.Name, repo.Created.Format(time.DateTime), repo.Deleted.Format(time.DateTime), repoid.String())
	if err != nil {
		log.
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
		return errors.ERR_BAD_DATA
	}
	row, err := d.db.Query("SELECT * FORM permitions WHERE id = ?", repoid)
	if err != nil {
		log.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error get repo")
		return repo, err
	}
	var identif string
	var id IDT
	err = row.Scan(&identif, &repo.Name, &repo.Created, &repo.Edited, &repo.Deleted)
	if err != nil {
		log.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error get repo")
		return repo, err
	}
	err = id.UnmarshalText([]byte(identif))
	repo.ID = id
	if err != nil {
		log.
			WithContext(ctx).
			WithField("repoid", repoid).
			WithError(err).
			Warn("error get repo")
		return repo, err
	}
	logrus.Trace("get info about Repositor ", repo.ID)
	return repo, nil
}

func (d *DB) GetRepos(ctx context.Context) ([]Repo, error) {
	var err error
	repos := make([]Repo, 0, 16)
	if err := ctx.Err(); err != nil {
		return repos, err
	}
	row, err := d.db.Query("SELECT * FORM permitions ")
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			Warn("error gets repo")
		return repos, err
	}
	for row.Next() {
		var repo Repo
		var identif string
		var id IDT
		err = row.Scan(&identif, &repo.Name, &repo.Created, &repo.Deleted)
		if err != nil {
			log.
				WithContext(ctx).
				WithError(err).
				Warn("error gets repo")
			return repos, err
		}
		err = id.UnmarshalText([]byte(identif))
		repo.ID = id
		if err != nil {
			log.WithContext(ctx).WithField("user_id", identif).WithError(err).Warn("error creating user")
			return repos, err
		}
		repos = append(repos, repo)
		logrus.Trace("get info about Repositor", repo.ID)
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
	statement, err := d.db.Prepare("INSERT INTO roles (role_id ,user_id, rep_id , branch, created, deleted) VALUES (?,?, ?,?, ?, ?)")
	if err != nil {
		log.
			WithContext(ctx).
			WithField("role ", ar.RoleID).
			WithError(err).
			Warn("error create role")
		return err
	}
	roleid := uuid.New().String()
	_, err = statement.Exec(roleid, ar.UserID, ar.RepoID, ar.Branches, ar.Created.Format(time.DateTime), ar.Deleted.Format(time.DateTime))
	if err != nil {
		log.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error create role")
		return err
	}
	logrus.Trace("Create Repositor", ar)
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
		log.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error edit role")
		return err
	}
	_, err = statement.Exec(ar.UserID.String(), ar.RepoID.String(), ar.Branches, ar.Created.Format(time.DateTime), ar.Deleted.Format(time.DateTime), ar.RoleID.String())
	if err != nil {
		log.
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
	statement, err := d.db.Prepare("DELETE * FORM roles WHERE role_id = ?")
	if err != nil {
		log.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error deleted role")
		return err
	}
	_, err = statement.Exec(roleid)
	if err != nil {
		log.
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
	var err error
	if err := ctx.Err(); err != nil {
		return role, err
	}
	if roleid == uuid.Nil {
		return errors.ERR_BAD_DATA
	}
	row := d.db.QueryRow("SELECT * FORM roles WHERE id = ?", roleid)
	if err != nil {
		log.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error get role")
		return role, err
	}
	var idrole, iduser, idrepo IDT
	var identifRole, identifUser, identifReposit string
	row.Scan(&identifRole, &identifUser, &identifReposit, &role.Branches, &role.Created, &role.Deleted)
	err = idrole.UnmarshalText([]byte(identifRole))
	if err != nil {
		log.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error get role")
		return role, err
	}
	err = iduser.UnmarshalText([]byte(identifUser))
	if err != nil {
		log.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error get role")
		return role, err
	}
	err = idrepo.UnmarshalText([]byte(identifReposit))
	if err != nil {
		log.
			WithContext(ctx).
			WithField("role ", roleid).
			WithError(err).
			Warn("error get role")
		return role, err
	}
	logrus.Trace("get info about Role ", role)
	return role, nil
}

func (d *DB) GetAccessRoles(ctx context.Context) ([]AccessRole, error) {
	var roles []AccessRole
	if err := ctx.Err(); err != nil {
		return roles, err
	}
	row, err := d.db.Query("SELECT * FORM roles ")
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			Warn("error get roles")
		return roles, err
	}
	for row.Next() {
		var idrole, iduser, idrepo IDT
		var identifRole, identifUser, identifReposit string
		var role AccessRole
		row.Scan(&identifRole, &identifUser, &identifReposit, &role.Branches, &role.Created, &role.Deleted)
		err = idrole.UnmarshalText([]byte(identifRole))
		if err != nil {
			log.
				WithContext(ctx).
				WithError(err).
				Warn("error get roles")
			return roles, err
		}
		err = iduser.UnmarshalText([]byte(identifUser))
		if err != nil {
			log.
				WithContext(ctx).
				WithError(err).
				Warn("error get roles")
			return roles, err
		}
		err = idrepo.UnmarshalText([]byte(identifReposit))
		if err != nil {
			log.
				WithContext(ctx).
				WithError(err).
				Warn("error get roles")
			return roles, err
		}
		if err != nil {
			log.WithContext(ctx).WithField("user_id", identifRole).WithError(err).Warn("error creating user")
			return roles, err
		}
		logrus.Trace("get info about Role ", role)
		roles = append(roles, role)
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
		return errors.ERR_BAD_DATA
	}
	row := d.db.QueryRow("SELECT userid.keys, users.name, users.email , users.created, users.edited, users.deleted FROM keys INNER JOIN users ON users.userid = keys.userid", key) // что делает и что выводит ?
	var answer string
	err = row.Scan(&answer, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
	if err != nil {
		log.
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
		return errors.ERR_BAD_DATA
	}
	row, err := d.db.Query("SELECT roleid FROM roles WHERE userid = ? AND WHERE repoid = ? AND WHERE branch = ?", userid.String(), repoid.String(), branch)
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithField("user_id ", userid).
			WithField("repo_id ", repoid).
			Warn("error check permitions")
		return false, err
	}
	var roleid string
	err = row.Scan(&roleid)
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithField("user_id ", userid).
			WithField("repo_id ", repoid).
			Warn("error check permitions")
		return false, err
	}
	logrus.Trace("get info permitions user ", userid)
	return true, err
}
