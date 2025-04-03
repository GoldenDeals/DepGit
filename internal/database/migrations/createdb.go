package migrations

type CreateDatabase interface {
	CreateTableUsers()
	CreateTableKeys()
	CreateTableRoles()
}

var log = logger.New("create db")

func (d *DB) CreateTableUsers() {
	statement, err := d.db.Prepare("IF NOT EXIST CREATE TABLE users (id VARCHAR(256), name VARCHAR(256), email VARCHAR(256), created TIME, edited TIME, deleted TIME)")
	if err != nil {
		log.Errorf("Error create table users ", err)
		return err
	}
	err = statement.Exec()
	if err != nil {
		log.Errorf("Error create table users ", err)
		return err
	}
}

func (d *DB) CreateTableKeys() {
	statement, err := d.db.Prepare("IF NOT EXIST CREATE TABLE keys (id VARCHAR(256),user_id VARCHAR(256), name VARCHAR(256), type INT, key VARCHAR(256),  created TIME, deleted TIME)")
	if err != nil {
		log.Errorf("Error create table keys ", err)
		return err
	}
	err = statement.Exec()
	if err != nil {
		log.Errorf("Error create table keys ", err)
		return err
	}
}

func (d *DB) CreateTableRoles() {
	statement, err := d.db.Prepare("IF NOT EXIST CREATE TABLE roles (role_id VARCHAR(256),user_id VARCHAR(256), rep_id VARCHAR(256), branch VARCHAR(256), created TIME, deleted TIME)")
	if err != nil {
		log.Errorf("Error create table roles ", err)
		return err
	}
	err = statement.Exec()
	if err != nil {
		log.Errorf("Error create table roles ", err)
		return err
	}
}
