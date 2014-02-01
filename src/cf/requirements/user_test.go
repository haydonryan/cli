package requirements

import (
	"cf"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testterm "testhelpers/terminal"
	"testing"
)

func TestUserReqExecute(t *testing.T) {
	user := cf.UserFields{}
	user.Username = "my-user"
	user.Guid = "my-user-guid"

	userRepo := &testapi.FakeUserRepository{FindByUsernameUserFields: user}
	ui := new(testterm.FakeUI)

	userReq := newUserRequirement("foo", ui, userRepo)
	success := userReq.Execute()

	assert.True(t, success)
	assert.Equal(t, userRepo.FindByUsernameUsername, "foo")
	assert.Equal(t, userReq.GetUser(), user)
}

func TestUserReqWhenUserDoesNotExist(t *testing.T) {
	userRepo := &testapi.FakeUserRepository{FindByUsernameNotFound: true}
	ui := new(testterm.FakeUI)

	testassert.AssertPanic(t, testterm.FailedWasCalled, func() {
		newUserRequirement("foo", ui, userRepo).Execute()
	})

	testassert.SliceContains(t, ui.Outputs, testassert.Lines{
		{"FAILED"},
		{"UserFields not found"},
	})
}
