package database_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	. "github.com/GoldenDeals/DepGit/internal/database"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	ase "github.com/stretchr/testify/assert"
)

var (
	log = logger.New("db_tests")
	d   *DB
)

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
	assert := ase.New(t)
	var ctx context.Context
	var err error
	user := NewUser("Jonson", "gayporno@yandex.ru")

	err = d.CreateUser(ctx, &user)
	assert.Nil(err)

	err = d.DeleteUser(ctx, user.ID)
	assert.Nil(err)
}

func TestsGetUser(t *testing.T) {
	assert := ase.New(t)
	var ctx context.Context
	var err error
	var checkuser User
	user := NewUser("Jonson", "gayporno@yandex.ru")

	err = d.CreateUser(ctx, &user)
	assert.Nil(err)

	checkuser, err = d.GetUser(ctx, user.ID)
	assert.Nil(err)
	assert.Equal(user, checkuser)
}

/*func TestsGetUsers(t *testing.T) {
	assert := ase.New(t)
	var ctx context.Context
	var err error
}
*/

func TestAddSshKey(t *testing.T) {
	assert := ase.New(t)
	var ctx context.Context
	var err error
	user := NewUser("Ahmed", "romalox@yandex.ru")
	data := make([]byte, 16)
	key := NewSShKey("lololoshka", RSA_SHA2_256, data)

	err = d.CreateUser(ctx, &user)
	assert.Nil(err)

	err = d.AddSshKey(ctx, user.ID, &key)
	assert.Nil(err)
}

func TestDeleteSshKey(t *testing.T) {
	assert := ase.New(t)
	var ctx context.Context
	var err error
	user := NewUser("Eblan", "bibobo@gmail.ru")
	data := make([]byte, 16)
	key := NewSShKey("lololoshka", RSA_SHA2_256, data)
	err = d.CreateUser(ctx, &user)
	assert.Nil(err)

	err = d.AddSshKey(ctx, user.ID, &key)
	assert.Nil(err)

	err = d.DeleteSshKey(ctx, key.ID)
	assert.Nil(err)
}

/*func TestGetSshKeys(t *testing.T){
	assert := ase.New(t)
	var ctx context.Context
	var err error
	user := NewUser("Ahmed", "romalox@yandex.ru")
	data := make([]byte, 16)
	key := NewSShKey("lololoshka", RSA_SHA2_256, data)

	err = d.CreateUser(ctx, &user)
	assert.Nil(err)

	err = d.AddSshKey(ctx, user.ID, &key)
	assert.Nil(err)

	keys := make(SshKey, 16)

}*/ // я не ебу как это тестить

func TestCreateRepo(t *testing.T) {
	assert := ase.New(t)
	var ctx context.Context
	var err error
	user := NewUser("Eblan", "bibobo@gmail.ru")
	data := make([]byte, 16)
	key := NewSShKey("lololoshka", RSA_SHA2_256, data)
	err = d.CreateUser(ctx, &user)
	assert.Nil(err)

	err = d.AddSshKey(ctx, user.ID, &key)
	assert.Nil(err)

	err = d.DeleteSshKey(ctx, key.ID)
	assert.Nil(err)
}
