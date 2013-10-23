package user_test

import (
	"cf"
	. "cf/commands/user"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func TestSetSpaceRoleFailsWithUsage(t *testing.T) {
	reqFactory, spaceRepo, userRepo := getSetSpaceRoleDeps()

	ui := callSetSpaceRole(t, []string{}, reqFactory, spaceRepo, userRepo)
	assert.True(t, ui.FailedWithUsage)

	ui = callSetSpaceRole(t, []string{"my-user"}, reqFactory, spaceRepo, userRepo)
	assert.True(t, ui.FailedWithUsage)

	ui = callSetSpaceRole(t, []string{"my-user", "my-org"}, reqFactory, spaceRepo, userRepo)
	assert.True(t, ui.FailedWithUsage)

	ui = callSetSpaceRole(t, []string{"my-user", "my-org", "my-space"}, reqFactory, spaceRepo, userRepo)
	assert.True(t, ui.FailedWithUsage)

	ui = callSetSpaceRole(t, []string{"my-user", "my-org", "my-space", "my-role"}, reqFactory, spaceRepo, userRepo)
	assert.False(t, ui.FailedWithUsage)
}

func TestSetSpaceRoleRequirements(t *testing.T) {
	args := []string{"username", "org", "space", "role"}
	reqFactory, spaceRepo, userRepo := getSetSpaceRoleDeps()

	reqFactory.LoginSuccess = false
	callSetSpaceRole(t, args, reqFactory, spaceRepo, userRepo)
	assert.False(t, testcmd.CommandDidPassRequirements)

	reqFactory.LoginSuccess = true
	callSetSpaceRole(t, args, reqFactory, spaceRepo, userRepo)
	assert.True(t, testcmd.CommandDidPassRequirements)

	assert.Equal(t, reqFactory.UserUsername, "username")
	assert.Equal(t, reqFactory.OrganizationName, "org")
}

func TestSetSpaceRole(t *testing.T) {
	args := []string{"some-user", "some-org", "some-space", "some-role"}
	reqFactory, spaceRepo, userRepo := getSetSpaceRoleDeps()

	reqFactory.LoginSuccess = true
	reqFactory.User = cf.User{Guid: "my-user-guid", Username: "my-user"}
	reqFactory.Organization = cf.Organization{Guid: "my-org-guid", Name: "my-org"}

	spaceRepo.FindByNameInOrgSpace = cf.Space{Guid: "my-space-guid", Name: "my-space"}

	ui := callSetSpaceRole(t, args, reqFactory, spaceRepo, userRepo)

	assert.Equal(t, spaceRepo.FindByNameInOrgName, "some-space")
	assert.Equal(t, spaceRepo.FindByNameInOrgOrg, reqFactory.Organization)

	assert.Contains(t, ui.Outputs[0], "Assigning role ")
	assert.Contains(t, ui.Outputs[0], "some-role")
	assert.Contains(t, ui.Outputs[0], "my-user")
	assert.Contains(t, ui.Outputs[0], "my-org")
	assert.Contains(t, ui.Outputs[0], "my-space")
	assert.Contains(t, ui.Outputs[0], "current-user")

	assert.Equal(t, userRepo.SetSpaceRoleUser, reqFactory.User)
	assert.Equal(t, userRepo.SetSpaceRoleSpace, spaceRepo.FindByNameInOrgSpace)
	assert.Equal(t, userRepo.SetSpaceRoleRole, "some-role")

	assert.Contains(t, ui.Outputs[1], "OK")
}

func getSetSpaceRoleDeps() (reqFactory *testreq.FakeReqFactory, spaceRepo *testapi.FakeSpaceRepository, userRepo *testapi.FakeUserRepository) {
	reqFactory = &testreq.FakeReqFactory{}
	spaceRepo = &testapi.FakeSpaceRepository{}
	userRepo = &testapi.FakeUserRepository{}
	return
}

func callSetSpaceRole(t *testing.T, args []string, reqFactory *testreq.FakeReqFactory, spaceRepo *testapi.FakeSpaceRepository, userRepo *testapi.FakeUserRepository) (ui *testterm.FakeUI) {
	ui = new(testterm.FakeUI)
	ctxt := testcmd.NewContext("set-space-role", args)

	token, err := testconfig.CreateAccessTokenWithTokenInfo(configuration.TokenInfo{
		Username: "current-user",
	})
	assert.NoError(t, err)

	config := &configuration.Configuration{
		Space:        cf.Space{Name: "my-space"},
		Organization: cf.Organization{Name: "my-org"},
		AccessToken:  token,
	}

	cmd := NewSetSpaceRole(ui, config, spaceRepo, userRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
