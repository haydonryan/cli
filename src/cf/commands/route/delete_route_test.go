package route_test

import (
	"cf"
	. "cf/commands/route"
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

func TestDeleteRouteRequirements(t *testing.T) {
	routeRepo := &testapi.FakeRouteRepository{}
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}

	callDeleteRoute(t, "y", []string{"-n", "my-host", "example.com"}, reqFactory, routeRepo)
	assert.True(t, testcmd.CommandDidPassRequirements)

	reqFactory = &testreq.FakeReqFactory{LoginSuccess: false}
	callDeleteRoute(t, "y", []string{"-n", "my-host", "example.com"}, reqFactory, routeRepo)
	assert.False(t, testcmd.CommandDidPassRequirements)
}

func TestDeleteRouteFailsWithUsage(t *testing.T) {
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
	routeRepo := &testapi.FakeRouteRepository{}
	ui := callDeleteRoute(t, "y", []string{}, reqFactory, routeRepo)
	assert.True(t, ui.FailedWithUsage)

	ui = callDeleteRoute(t, "y", []string{"example.com"}, reqFactory, routeRepo)
	assert.False(t, ui.FailedWithUsage)

	ui = callDeleteRoute(t, "y", []string{"-n", "my-host", "example.com"}, reqFactory, routeRepo)
	assert.False(t, ui.FailedWithUsage)
}

func TestDeleteRouteWithConfirmation(t *testing.T) {
	domain := cf.DomainFields{}
	domain.Guid = "domain-guid"
	domain.Name = "example.com"
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
	route := cf.Route{}
	route.Guid = "route-guid"
	route.Host = "my-host"
	route.Domain = domain
	routeRepo := &testapi.FakeRouteRepository{
		FindByHostAndDomainRoute: route,
	}

	ui := callDeleteRoute(t, "y", []string{"-n", "my-host", "example.com"}, reqFactory, routeRepo)

	testassert.SliceContains(t, ui.Prompts, testassert.Lines{
		{"Really delete", "my-host"},
	})

	testassert.SliceContains(t, ui.Outputs, testassert.Lines{
		{"Deleting route", "my-host.example.com"},
		{"OK"},
	})
	assert.Equal(t, routeRepo.DeleteRouteGuid, "route-guid")
}

func TestDeleteRouteWithForce(t *testing.T) {
	domain := cf.DomainFields{}
	domain.Guid = "domain-guid"
	domain.Name = "example.com"
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
	route := cf.Route{}
	route.Guid = "route-guid"
	route.Host = "my-host"
	route.Domain = domain
	routeRepo := &testapi.FakeRouteRepository{
		FindByHostAndDomainRoute: route,
	}

	ui := callDeleteRoute(t, "", []string{"-f", "-n", "my-host", "example.com"}, reqFactory, routeRepo)

	assert.Equal(t, len(ui.Prompts), 0)

	testassert.SliceContains(t, ui.Outputs, testassert.Lines{
		{"Deleting", "my-host.example.com"},
		{"OK"},
	})
	assert.Equal(t, routeRepo.DeleteRouteGuid, "route-guid")
}

func TestDeleteRouteWhenRouteDoesNotExist(t *testing.T) {
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
	routeRepo := &testapi.FakeRouteRepository{
		FindByHostAndDomainNotFound: true,
	}

	ui := callDeleteRoute(t, "y", []string{"-n", "my-host", "example.com"}, reqFactory, routeRepo)

	testassert.SliceContains(t, ui.Outputs, testassert.Lines{
		{"Deleting", "my-host.example.com"},
		{"OK"},
		{"my-host", "does not exist"},
	})
}

func callDeleteRoute(t *testing.T, confirmation string, args []string, reqFactory *testreq.FakeReqFactory, routeRepo *testapi.FakeRouteRepository) (ui *testterm.FakeUI) {
	ui = &testterm.FakeUI{
		Inputs: []string{confirmation},
	}
	ctxt := testcmd.NewContext("delete-route", args)

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

	cmd := NewDeleteRoute(ui, config, routeRepo)

	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
