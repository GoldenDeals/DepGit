package database

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

type TestingAccessMetods interface {
	TestsCreateUser(t *testing.T)
	TestsEditUser(t *testing.T)
	TestsDeleteUser(t *testing.T)
	TestsGetUser(t *testing.T)
	TestsGetUsers(t *testing.T)

	TestsAddSshKey(t *testing.T)
	TestsDeleteSshKey(t *testing.T)
	TestsGetSshKeys(t *testing.T)

	TestsCreateRepo(t *testing.T)
	TestsDeleteRepo(t *testing.T)
	TestsUpdateRepo(t *testing.T)
	TestsGetRepo(t *testing.T)
	TestsGetRepos(t *testing.T)

	TestsCreateAccessRole(t *testing.T)
	TestsEditAccessRole(t *testing.T)
	TestsDeleteAccessRole(t *testing.T)
	TestsGetAccessRole(t *testing.T)
	TestsGetAccessRoles(t *testing.T)

	TestsUserByKey(t *testing.T)
	TestsCheckPermissions(t *testing.T)
}

var log = logger.New("db_tests")

func (d *DB) TestsCreateUser(t *testing.T) {
	var user User
	var ctx context.Context
	var err error
	user.Create("Jonson", "gayporno@yandex.ru")
	if err != nil {
		log.Errorf("Error  ", err)
	}
	err = d.CreateUser(ctx, &user)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	checkuser, err := d.GetUser(ctx, user.ID)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	if user != checkuser {
		log.Errorf("Error  ", err)
	}
	logrus.Trace("Tests Create User Complete")
}

func (d *DB) TestsEditUser(t *testing.T) {
	var user, checkuser User
	var ctx context.Context
	var err error
	user.Create("Jonson", "gayporno@yandex.ru")
	if err != nil {
		log.Errorf("Error  ", err)
	}
	checkuser.Create("Soprano", "gayporno@gmail.ru")
	if err != nil {
		log.Errorf("Error  ", err)
	}
	err = d.CreateUser(ctx, &user)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	err = d.EditUser(ctx, user.ID, &checkuser)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	checkwithuser, err := d.GetUser(ctx, user.ID)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	if checkwithuser != checkuser {
		log.Errorf("Error  ", err)
	}
	logrus.Trace("Tests Edit User Complete")
}

func (d *DB) TestsDeleteUser(t *testing.T) {
	var user User
	var ctx context.Context
	var err error
	var number int
	user.Create("Jonson", "gayporno@yandex.ru")
	if err != nil {
		log.Errorf("Error  ", err)
	}
	err = d.CreateUser(ctx, &user)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	err = d.DeleteUser(ctx, user.ID)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	row := d.db.QueryRow("SELECT COUNT(userid) FROM users WHERE userid = ?", user.ID)
	err = row.Scan(&number)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	if number != 0 {
		log.Errorf("Error  ", err)
	}
	logrus.Trace("Tests Delete User Complete")
}

func (d *DB) TestsGetUser(t *testing.T) {
	var user User
	var ctx context.Context
	var err error
	user.Create("Jonson", "gayporno@yandex.ru")
	if err != nil {
		log.Errorf("Error  ", err)
	}
	err = d.CreateUser(ctx, &user)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	checkuser, err := d.GetUser(ctx, user.ID)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	if user != checkuser {
		log.Errorf("Error  ", err)
	}
	logrus.Trace("Tests Get User Complete")
}

func (d *DB) TestsGetUsers(t *testing.T) {
	var user User
	var ctx context.Context
	var err error
	user.Create("Jonson", "gayporno@yandex.ru")
	if err != nil {
		log.Errorf("Error  ", err)
	}
	err = d.CreateUser(ctx, &user)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	checkuser, err := d.GetUser(ctx, user.ID)
	if err != nil {
		log.Errorf("Error  ", err)
	}
	if user != checkuser {
		log.Errorf("Error  ", err)
	}
	logrus.Trace("Tests Create User Complete")
}
