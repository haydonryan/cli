package organization_test

import (
	"cf"
	"cf/commands/organization"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func TestListQuotasRequirements(t *testing.T) {
	quotaRepo := &testapi.FakeQuotaRepository{}

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
	callListQuotas(t, reqFactory, quotaRepo)
	assert.True(t, testcmd.CommandDidPassRequirements)

	reqFactory = &testreq.FakeReqFactory{LoginSuccess: false}
	callListQuotas(t, reqFactory, quotaRepo)
	assert.False(t, testcmd.CommandDidPassRequirements)
}

func TestListQuotas(t *testing.T) {
	quota := cf.QuotaFields{}
	quota.Name = "quota-name"
	quota.MemoryLimit = 1024

	quotaRepo := &testapi.FakeQuotaRepository{FindAllQuotas: []cf.QuotaFields{quota}}
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
	ui := callListQuotas(t, reqFactory, quotaRepo)

	testassert.SliceContains(t, ui.Outputs, testassert.Lines{
		{"Getting quotas as", "my-user"},
		{"OK"},
		{"name", "memory limit"},
		{"quota-name", "1g"},
	})
}

func callListQuotas(t *testing.T, reqFactory *testreq.FakeReqFactory, quotaRepo *testapi.FakeQuotaRepository) (fakeUI *testterm.FakeUI) {
	fakeUI = &testterm.FakeUI{}
	ctxt := testcmd.NewContext("quotas", []string{})

	token, err := testconfig.CreateAccessTokenWithTokenInfo(configuration.TokenInfo{
		Username: "my-user",
	})
	assert.NoError(t, err)

	spaceFields := cf.SpaceFields{}
	spaceFields.Name = "my-space"

	orgFields := cf.OrganizationFields{}
	orgFields.Name = "my-org"

	config := &configuration.Configuration{
		SpaceFields:        spaceFields,
		OrganizationFields: orgFields,
		AccessToken:        token,
	}

	cmd := organization.NewListQuotas(fakeUI, config, quotaRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
