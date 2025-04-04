package database

import "testing"

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
