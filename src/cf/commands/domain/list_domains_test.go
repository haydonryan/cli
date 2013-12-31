package domain_test

import (
	"cf"
	"cf/commands/domain"
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

func TestListDomainsRequirements(t *testing.T) {
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: true}
	domainRepo := &testapi.FakeDomainRepository{}

	callListDomains(t, []string{}, reqFactory, domainRepo)
	assert.True(t, testcmd.CommandDidPassRequirements)

	reqFactory = &testreq.FakeReqFactory{LoginSuccess: false, TargetedOrgSuccess: true}
	callListDomains(t, []string{}, reqFactory, domainRepo)
	assert.False(t, testcmd.CommandDidPassRequirements)

	reqFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: false}
	callListDomains(t, []string{}, reqFactory, domainRepo)
	assert.False(t, testcmd.CommandDidPassRequirements)
}

func TestListDomainsFailsWithUsage(t *testing.T) {
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: true}
	domainRepo := &testapi.FakeDomainRepository{}

	ui := callListDomains(t, []string{"foo"}, reqFactory, domainRepo)
	assert.True(t, ui.FailedWithUsage)
}

func TestListDomains(t *testing.T) {
	orgFields := cf.OrganizationFields{}
	orgFields.Name = "my-org"
	orgFields.Guid = "my-org-guid"

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: true, OrganizationFields: orgFields}
	domain1 := cf.Domain{}
	domain1.Shared = true
	domain1.Name = "Domain1"

	domain2 := cf.Domain{}
	domain2.Shared = false
	domain2.Name = "Domain2"

	space1 := cf.SpaceFields{}
	space1.Name = "my-space"

	space2 := cf.SpaceFields{}
	space2.Name = "my-space-2"

	domain2.Spaces = []cf.SpaceFields{space1, space2}

	domain3 := cf.Domain{}
	domain3.Shared = false
	domain3.Name = "Domain3"

	domainRepo := &testapi.FakeDomainRepository{
		ListDomainsForOrgDomains: []cf.Domain{domain1, domain2, domain3},
	}
	ui := callListDomains(t, []string{}, reqFactory, domainRepo)

	assert.Equal(t, domainRepo.ListDomainsForOrgDomainsGuid, "my-org-guid")
	testassert.SliceContains(t, ui.Outputs, testassert.Lines{
		{"Getting domains in org", "my-org", "my-user"},
		{"name", "status", "spaces"},
		{"Domain1", "shared"},
		{"Domain2", "owned", "my-space", "my-space-2"},
		{"Domain3", "reserved"},
	})
}

func callListDomains(t *testing.T, args []string, reqFactory *testreq.FakeReqFactory, domainRepo *testapi.FakeDomainRepository) (fakeUI *testterm.FakeUI) {
	fakeUI = new(testterm.FakeUI)
	ctxt := testcmd.NewContext("domains", args)

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

	cmd := domain.NewListDomains(fakeUI, config, domainRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
