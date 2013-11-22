package service_test

import (
	"cf"
	"cf/api"
	. "cf/commands/service"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func TestCreateUserProvidedServiceWithParameterList(t *testing.T) {
	repo := &testapi.FakeUserProvidedServiceInstanceRepo{}
	fakeUI := callCreateUserProvidedService(t,
		[]string{"-p", `"foo, bar, baz"`, "my-custom-service"},
		[]string{"foo value", "bar value", "baz value"},
		repo,
	)

	assert.Contains(t, fakeUI.Prompts[0], "foo")
	assert.Contains(t, fakeUI.Prompts[1], "bar")
	assert.Contains(t, fakeUI.Prompts[2], "baz")

	assert.Equal(t, repo.CreateName, "my-custom-service")
	assert.Equal(t, repo.CreateParams, map[string]string{
		"foo": "foo value",
		"bar": "bar value",
		"baz": "baz value",
	})

	assert.Contains(t, fakeUI.Outputs[0], "Creating user provided service")
	assert.Contains(t, fakeUI.Outputs[0], "my-custom-service")
	assert.Contains(t, fakeUI.Outputs[0], "my-org")
	assert.Contains(t, fakeUI.Outputs[0], "my-space")
	assert.Contains(t, fakeUI.Outputs[0], "my-user")
	assert.Contains(t, fakeUI.Outputs[1], "OK")
}

func TestCreateUserProvidedServiceWithJson(t *testing.T) {
	repo := &testapi.FakeUserProvidedServiceInstanceRepo{}
	fakeUI := callCreateUserProvidedService(t,
		[]string{"-p", `{"foo": "foo value", "bar": "bar value", "baz": "baz value"}`, "my-custom-service"},
		[]string{},
		repo,
	)

	assert.Empty(t, fakeUI.Prompts)

	assert.Equal(t, repo.CreateName, "my-custom-service")
	assert.Equal(t, repo.CreateParams, map[string]string{
		"foo": "foo value",
		"bar": "bar value",
		"baz": "baz value",
	})

	assert.Contains(t, fakeUI.Outputs[0], "Creating user provided service")
	assert.Contains(t, fakeUI.Outputs[1], "OK")
}

func TestCreateUserProvidedServiceWithNoSecondArgument(t *testing.T) {
	userProvidedServiceInstanceRepo := &testapi.FakeUserProvidedServiceInstanceRepo{}
	fakeUI := callCreateUserProvidedService(t,
		[]string{"my-custom-service"},
		[]string{},
		userProvidedServiceInstanceRepo,
	)

	assert.Contains(t, fakeUI.Outputs[0], "Creating user provided service")
	assert.Contains(t, fakeUI.Outputs[1], "OK")
}

func TestCreateUserProvidedServiceWithSyslogDrain(t *testing.T) {
	repo := &testapi.FakeUserProvidedServiceInstanceRepo{}

	fakeUI := callCreateUserProvidedService(t,
		[]string{"-l", "syslog://example.com", "-p", `{"foo": "foo value", "bar": "bar value", "baz": "baz value"}`, "my-custom-service"},
		[]string{},
		repo,
	)
	assert.Equal(t, repo.CreateDrainUrl, "syslog://example.com")
	assert.Contains(t, fakeUI.Outputs[0], "Creating user provided service")
	assert.Contains(t, fakeUI.Outputs[1], "OK")
}

func callCreateUserProvidedService(t *testing.T, args []string, inputs []string, userProvidedServiceInstanceRepo api.UserProvidedServiceInstanceRepository) (fakeUI *testterm.FakeUI) {
	fakeUI = &testterm.FakeUI{Inputs: inputs}
	ctxt := testcmd.NewContext("create-user-provided-service", args)
	reqFactory := &testreq.FakeReqFactory{}

	token, err := testconfig.CreateAccessTokenWithTokenInfo(configuration.TokenInfo{
		Username: "my-user",
	})
	assert.NoError(t, err)
	org := cf.OrganizationFields{}
	org.Name = "my-org"
	space := cf.SpaceFields{}
	space.Name = "my-space"
	config := &configuration.Configuration{
		SpaceFields:        space,
		OrganizationFields: org,
		AccessToken:        token,
	}

	cmd := NewCreateUserProvidedService(fakeUI, config, userProvidedServiceInstanceRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
