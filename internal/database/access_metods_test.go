package database_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	. "github.com/GoldenDeals/DepGit/internal/database"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/sirupsen/logrus"
	ase "github.com/stretchr/testify/assert"
)

var log = logger.New("db_tests")
var d *DB

func TestMain(m *testing.M) {
	err := d.Init("./db.sqllite3")
	if err != nil {
		fmt.Print(err)
		os.Exit(255)
	}
	code := m.Run()
	err = d.Close()
	if err != nil {
		fmt.Print(err)
	}

	os.Exit(code)
}

func TestsCreateUser(t *testing.T) {
	assert := ase.New(t)
	var ctx context.Context
	var err error
	user := NewUser("Jonson", "gayporno@yandex.ru")
	err = d.CreateUser(ctx, &user)
	assert.Nil(err)

	checkuser, err := d.GetUser(ctx, user.ID)
	assert.Nil(err)
	assert.Equal(user, checkuser)
}

func TestsEditUser(t *testing.T) {
	assert := ase.New(t)

	var ctx context.Context
	var err error
	user := NewUser("Jonson", "gayporno@yandex.ru")
	checkuser := NewUser("Soprano", "gayporno@gmail.ru")

	err = d.CreateUser(ctx, &user)
	assert.Nil(err)

	err = d.EditUser(ctx, user.ID, &checkuser)
	assert.Nil(err)

	checkwithuser, err := d.GetUser(ctx, user.ID)
	assert.Nil(err)
	assert.Equal(checkwithuser, checkuser)
}

func TestsDeleteUser(t *testing.T) {
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

	d.GetUser(context.Background(), user.ID)
	if number != 0 {
		log.Errorf("Error  ", err)
	}
	logrus.Trace("Tests Delete User Complete")
}

func TestsGetUser(t *testing.T) {
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

func TestsGetUsers(t *testing.T) {
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
