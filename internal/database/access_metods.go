package database

import (
	"database/sql"
	"time"

	"github.com/gobwas/glob"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	SSH_KEY_TYPE_RSA    SSH_KEY_TYPE = iota + 1
	RSA_SHA2_256        SSH_KEY_TYPE
	RSA_SHA2_512        SSH_KEY_TYPE
	SSH_RSA             SSH_KEY_TYPE
	ECDSA_SHA2_NISTP256 SSH_KEY_TYPE
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

type AccessManager interface {
	UserByKey(key []byte) (*User, error)
	CheckPermissions(userid IDT, repoid IDT, branch string) (bool, error)
}

type DB struct {
	db *sql.DB
}

func (d *DB) Init() error {
	var err error
	d.db, err = sql.Open("sqlite3", "./migrations/access.db")
	if err != nil {
		log.Errorf("Error init Database sql open ", err)
		return err
	}
	err = d.db.Ping()
	if err != nil {
		log.Errorf("Error init Database ping ", err)
		return err
	}
	logrus.Debug("Database access.db open")
	return nil
}

func (d *DB) CreateUser(user *User) error {
	statement, err := d.db.Prepare("INSERT INTO users (id, name, email, created, edited, deleted) VALUES (?,?,?,?,?)")
	if err != nil {
		log.Errorf("Error Create User ", user.ID, err)
		return err
	}
	userid := uuid.New().String()
	_, err = statement.Exec(userid, user.Name, time.Now().Format(time.DateTime), user.Edited.Format(time.DateTime), user.Deleted.Format(time.DateTime))
	if err != nil {
		log.Errorf("Error Execute sql request ", err)
		return err
	}
	logrus.Trace("Create User", userid, user.Name)
	return nil
}

func (d *DB) EditUser(userid IDT, user *User) error {
	statement, err := d.db.Prepare("UPDATE users SET name = ? , email = ? , created = ?, edited = ?, deleted = ? WHERE id = ?")
	if err != nil {
		log.Errorf("Error edit user ", user.ID, err)
		return err
	}
	_, err = statement.Exec(user.Name, user.Email, user.Created.Format(time.DateTime), time.Now().Format(time.DateTime), user.Deleted.Format(time.DateTime), userid.String())
	if err != nil {
		log.Errorf("Error Execute sql request ", err)
		return err
	}
	logrus.Trace("Edit user", user.ID, user.Name)
	return nil
}

func (d *DB) DeleteUser(userid IDT) error {
	statement, err := d.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		log.Errorf("Error delete User ", userid, err)
		return err
	}
	_, err = statement.Exec(userid)
	if err != nil {
		log.Errorf("Error Execute sql request ", err)
		return err
	}
	logrus.Trace("Delete user", userid)
	return nil
}

func (d *DB) GetUser(userid IDT) (User, error) {
	var user User
	var err error
	var identif string
	var id IDT
	row := d.db.QueryRow("SELECT * FORM users WHERE id = ?", userid)
	err = row.Scan(&identif, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
	if err != nil {
		log.Errorf("Error Scan sql request ", err)
		return user, err
	}
	err = id.UnmarshalText([]byte(identif))
	user.ID = id
	if err != nil {
		log.Errorf("Error Get user", userid, err)
		return user, err
	}
	logrus.Trace("get info about user", user)
	return user, nil
}

func (d *DB) GetUsers() ([]User, error) {
	users := make([]User, 0, 16)
	row, err := d.db.Query("SELECT * FORM users ")
	if err != nil {
		log.Errorf("Error GetUsers ", err)
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
			log.Errorf("Error scan GetUsers", identif, err)
			return users, err
		}
		users = append(users, user)
		logrus.Trace("get info about user", user.ID, user.Name)
	}

	return users, nil
}

func (d *DB) AddSshKey(userid IDT, key *SshKey) error {
	var err error
	statement, err := d.db.Prepare("INSERT INTO keys (id,user_id, name, type, key,  created, deleted) VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.Errorf("Error AddSshKey ", userid, key, err)
		return err
	}
	_, err = statement.Exec(key.ID.String(), key.UserID.String(), key.Name, key.Type, key.Data, time.Now().Format(time.DateTime), key.Deleted.Format(time.DateTime))
	if err != nil {
		log.Errorf("Error Execute sql request ", err) //
		return err
	}
	logrus.Trace("Add SSH KEY", key, userid)
	return nil
}

func (d *DB) DelSshKey(userid IDT, keyid IDT) error {
	statement, err := d.db.Prepare("DELETE * FORM keys WHERE id = ?, AND WHERE user_id = ?")
	if err != nil {
		log.Errorf("Error DelSshKey ", userid, keyid, err)
		return err
	}
	_, err = statement.Exec(keyid, userid)
	if err != nil {
		log.Errorf("Error Execute sql request ", err) //
		return err
	}
	logrus.Trace("Delete SSH key ", userid, keyid)
	return nil
}

func (d *DB) GetSshKeys(userid IDT) ([]SshKey, error) {
	keys := make([]SshKey, 0, 16)
	row, err := d.db.Query("SELECT * FROM keys WHERE user_id = ?", userid)
	if err != nil {
		log.Errorf("Error GetSshKeys ", userid, err)
		return keys, err
	}
	for row.Next() {
		var key SshKey
		err = row.Scan(&key.ID, &key.UserID, &key.Name, &key.Type, &key.Data, &key.Created, &key.Deleted)
		if err != nil {
			log.Errorf("Error Scan sql request ", err) //
			return keys, err
		}
		keys = append(keys, key)
	}
	logrus.Trace("Get SSH keys user ", userid)
	return keys, nil
}

func (d *DB) CreateRepo(repo *Repo) error {
	statement, err := d.db.Prepare("INSERT INTO permitions (id, name, created, edited ,deleted) VALUES (?, ?,?, ?, ?)")
	if err != nil {
		log.Errorf("Error Create Repo ", repo.ID, err)
		return err
	}
	repoid := uuid.New().String()
	_, err = statement.Exec(repoid, repo.Name, time.Now().Format(time.DateTime), repo.Edited, repo.Deleted.Format(time.DateTime))
	if err != nil {
		log.Errorf("Error Execute sql request ", err) //
		return err
	}
	logrus.Trace("Create Repositor", repo.ID)
	return nil
}

func (d *DB) DeleteRepo(repoid IDT) error {
	statement, err := d.db.Prepare("DELETE * FORM permitions WHERE id = ?")
	if err != nil {
		log.Errorf("Error DeleteRepo ", repoid, err)
		return err
	}
	_, err = statement.Exec(repoid)
	if err != nil {
		log.Errorf("Error Execute sql request ", err) //
		return err
	}
	logrus.Trace("Deleted.Format(time.DateTime) Repositor", repoid)
	return nil
}

func (d *DB) UpdateRepo(repoid IDT, repo Repo) error {
	statement, err := d.db.Prepare("UPDATE permitions SET name = ? , created = ?, deleted = ? WHERE id = ?")
	if err != nil {
		log.Errorf("Error UpdateRepo ", repoid, err)
		return err
	}
	_, err = statement.Exec(repo.Name, repo.Created.Format(time.DateTime), repo.Deleted.Format(time.DateTime), repoid.String())
	if err != nil {
		log.Errorf("Error Execute sql request ", err) //
		return err
	}
	logrus.Trace("Edit Repositor ", repo.ID)
	return nil
}

func (d *DB) GetRepo(repoid IDT) (Repo, error) {
	var repo Repo
	row, err := d.db.Query("SELECT * FORM permitions WHERE id = ?", repoid)
	if err != nil {
		log.Errorf("Error GetRepo ", repoid, err)
		return repo, err
	}
	var identif string
	var id IDT
	err = row.Scan(&identif, &repo.Name, &repo.Created, &repo.Edited, &repo.Deleted)
	if err != nil {
		log.Errorf("Error scan sql request ", err) //
		return repo, err
	}
	err = id.UnmarshalText([]byte(identif))
	repo.ID = id
	if err != nil {
		log.Errorf("Error GetRepo scan", identif, repoid, err)
		return repo, err
	}
	logrus.Trace("get info about Repositor ", repo.ID)
	return repo, nil
}

func (d *DB) GetRepos() ([]Repo, error) {
	var err error
	repos := make([]Repo, 0, 16)
	row, err := d.db.Query("SELECT * FORM permitions ")
	if err != nil {
		log.Errorf("Error GetRepos ", err)
		return repos, err
	}
	for row.Next() {
		var repo Repo
		var identif string
		var id IDT
		err = row.Scan(&identif, &repo.Name, &repo.Created, &repo.Deleted)
		if err != nil {
			log.Errorf("Error scan sql request ", err) //
			return repos, err
		}
		err = id.UnmarshalText([]byte(identif))
		repo.ID = id
		if err != nil {
			log.Errorf("Error GetRepos scan ", identif, err)
			return repos, err
		}
		repos = append(repos, repo)
		logrus.Trace("get info about Repositor", repo.ID)
	}

	return repos, nil
}

func (d *DB) CreateAccessRole(ar *AccessRole) error {
	statement, err := d.db.Prepare("INSERT INTO roles (role_id ,user_id, rep_id , branch, created, deleted) VALUES (?,?, ?,?, ?, ?)")
	if err != nil {
		log.Errorf("Error CreateAccessRole ", ar.RepoID, ar.RepoID, ar.UserID, err)
		return err
	}
	roleid := uuid.New().String()
	_, err = statement.Exec(roleid, ar.UserID, ar.RepoID, ar.Branches, ar.Created.Format(time.DateTime), ar.Deleted.Format(time.DateTime))
	if err != nil {
		log.Errorf("Error Execute sql request ", err) //
		return err
	}
	logrus.Trace("Create Repositor", ar)
	return nil
}

func (d *DB) EditAccessRole(roleid IDT, ar *AccessRole) error {
	statement, err := d.db.Prepare("UPDATE roles SET user_id = ?, rep_id = ?, branch = ?, created = ?, deleted = ? WHERE role_id = ?")
	if err != nil {
		log.Errorf("Error EditAccessRole ", roleid, err)
		return err
	}
	_, err = statement.Exec(ar.UserID.String(), ar.RepoID.String(), ar.Branches, ar.Created.Format(time.DateTime), ar.Deleted.Format(time.DateTime), ar.RoleID.String())
	if err != nil {
		log.Errorf("Error Execute sql request ", err) //
		return err
	}
	logrus.Trace("Edit role", ar.RoleID)
	return nil
}

func (d *DB) DeleteAccessRole(roleid IDT) error {
	statement, err := d.db.Prepare("DELETE * FORM roles WHERE role_id = ?")
	if err != nil {
		log.Errorf("Error DeleteAccessRole ", roleid, err)
		return err
	}
	_, err = statement.Exec(roleid)
	if err != nil {
		log.Errorf("Error Execute sql request ", err) //
		return err
	}
	logrus.Trace("Deleted.Format(time.DateTime) Role", roleid)
	return nil
}

func (d *DB) GetAccessRole(roleid IDT) (AccessRole, error) {
	var role AccessRole
	var err error
	row := d.db.QueryRow("SELECT * FORM roles WHERE id = ?", roleid)
	if err != nil {
		log.Errorf("Error GetAccessRoles ", roleid, err)
		return role, err
	}
	var idrole, iduser, idrepo IDT
	var identifRole, identifUser, identifReposit string
	row.Scan(&identifRole, &identifUser, &identifReposit, &role.Branches, &role.Created, &role.Deleted)
	err = idrole.UnmarshalText([]byte(identifRole))
	if err != nil {
		log.Errorf("Error UnmarshalText sql request ", err) //
		return role, err
	}
	err = iduser.UnmarshalText([]byte(identifUser))
	if err != nil {
		log.Errorf("Error UnmarshalText sql request ", err) //
		return role, err
	}
	err = idrepo.UnmarshalText([]byte(identifReposit))
	if err != nil {
		log.Errorf("Error GetAccessRoles scan ", identifRole, err)
		return role, err
	}
	logrus.Trace("get info about Role ", role)
	return role, nil
}

func (d *DB) GetAccessRoles() ([]AccessRole, error) {
	var roles []AccessRole
	row, err := d.db.Query("SELECT * FORM roles ")
	if err != nil {
		log.Errorf("Error GetAccessRoles ", err)
		return roles, err
	}
	for row.Next() {
		var idrole, iduser, idrepo IDT
		var identifRole, identifUser, identifReposit string
		var role AccessRole
		row.Scan(&identifRole, &identifUser, &identifReposit, &role.Branches, &role.Created, &role.Deleted)
		err = idrole.UnmarshalText([]byte(identifRole))
		if err != nil {
			log.Errorf("Error UnmarshalText sql request ", err) //
			return roles, err
		}
		err = iduser.UnmarshalText([]byte(identifUser))
		if err != nil {
			log.Errorf("Error UnmarshalText sql request ", err) //
			return roles, err
		}
		err = idrepo.UnmarshalText([]byte(identifReposit))
		if err != nil {
			log.Errorf("Error UnmarshalText sql request ", err) //
			return roles, err
		}
		if err != nil {
			log.Errorf("Error GetAccessRoles scan ", identifRole, err)
			return roles, err
		}
		logrus.Trace("get info about Role ", role)
		roles = append(roles, role)
	}
	return roles, nil
}

func (d *DB) UserByKey(key []byte) (User, error) {
	var user User
	var err error
	row := d.db.QueryRow("SELECT userid.keys, users.name, users.email , users.created, users.edited, users.deleted FROM keys INNER JOIN users ON users.userid = keys.userid", key) // что делает и что выводит ?
	var answer string
	err = row.Scan(&answer, &user.Name, &user.Email, &user.Created, &user.Edited, &user.Deleted)
	if err != nil {
		log.Errorf("Error Scan sql request ", key, err)
		return user, err
	}
	return user, nil
}

func (d *DB) CheckPermissions(userid IDT, repoid IDT, branch string) (bool, error) {
	row, err := d.db.Query("SELECT roleid FROM roles WHERE userid = ? AND WHERE repoid = ? AND WHERE branch = ?", userid.String(), repoid.String(), branch)
	if err != nil {
		log.Errorf("Error CheckPermissions ", userid, repoid, branch, err)
		return false, err
	}
	var roleid string
	err = row.Scan(&roleid)
	if err != nil {
		log.Errorf("Error CheckPermissions scan", userid, repoid, branch, roleid, err)
		return false, err
	}
	return true, err
}
