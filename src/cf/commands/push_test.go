package commands_test

import (
	"cf"
	"cf/api"
	. "cf/commands"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	"testhelpers"
	"testing"
)

func TestPushingAppWhenItDoesNotExist(t *testing.T) {
	domains := []cf.Domain{
		cf.Domain{Name: "foo.cf-app.com", Guid: "foo-domain-guid"},
	}
	domainRepo := &testhelpers.FakeDomainRepository{FindByNameDomain: domains[0]}
	routeRepo := &testhelpers.FakeRouteRepository{}
	appRepo := &testhelpers.FakeApplicationRepository{AppByNameErr: true}

	fakeUI := callPush([]string{"--name", "my-new-app"}, basePushConfig(), appRepo, domainRepo, routeRepo)

	assert.Contains(t, fakeUI.Outputs[0], "Creating my-new-app...")
	assert.Equal(t, appRepo.CreatedApp.Name, "my-new-app")
	assert.Equal(t, appRepo.CreatedApp.Instances, 1)
	assert.Contains(t, fakeUI.Outputs[1], "OK")

	assert.Contains(t, fakeUI.Outputs[2], "Creating route my-new-app.foo.cf-app.com...")
	assert.Equal(t, routeRepo.CreatedRoute.Host, "my-new-app")
	assert.Equal(t, routeRepo.CreatedRouteDomain.Guid, "foo-domain-guid")
	assert.Contains(t, fakeUI.Outputs[3], "OK")

	assert.Contains(t, fakeUI.Outputs[4], "Binding my-new-app.foo.cf-app.com to my-new-app...")
	assert.Equal(t, routeRepo.BoundApp.Name, "my-new-app")
	assert.Equal(t, routeRepo.BoundRoute.Host, "my-new-app")
	assert.Contains(t, fakeUI.Outputs[5], "OK")

	assert.Contains(t, fakeUI.Outputs[6], "Uploading my-new-app...")
	assert.Equal(t, appRepo.UploadedApp.Guid, "my-new-app-guid")
	assert.Contains(t, fakeUI.Outputs[7], "OK")
}

func TestPushingAppWithCustomFlags(t *testing.T) {
	domain := cf.Domain{Name: "bar.cf-app.com", Guid: "bar-domain-guid"}
	domainRepo := &testhelpers.FakeDomainRepository{FindByNameDomain: domain}
	routeRepo := &testhelpers.FakeRouteRepository{}
	appRepo := &testhelpers.FakeApplicationRepository{AppByNameErr: true}

	fakeUI := callPush([]string{
		"--name", "my-new-app",
		"--domain", "bar.cf-app.com",
		"--host", "my-hostname",
		"--instances", "3",
	}, basePushConfig(), appRepo, domainRepo, routeRepo)

	assert.Contains(t, fakeUI.Outputs[0], "Creating my-new-app...")
	assert.Equal(t, appRepo.CreatedApp.Name, "my-new-app")
	assert.Equal(t, appRepo.CreatedApp.Instances, 3)
	assert.Contains(t, fakeUI.Outputs[1], "OK")

	assert.Contains(t, fakeUI.Outputs[2], "Creating route my-hostname.bar.cf-app.com...")
	assert.Equal(t, domainRepo.FindByNameName, "bar.cf-app.com")
	assert.Equal(t, routeRepo.CreatedRoute.Host, "my-hostname")
	assert.Equal(t, routeRepo.CreatedRouteDomain.Guid, "bar-domain-guid")
	assert.Contains(t, fakeUI.Outputs[3], "OK")

	assert.Contains(t, fakeUI.Outputs[4], "Binding my-hostname.bar.cf-app.com to my-new-app...")
	assert.Equal(t, routeRepo.BoundApp.Name, "my-new-app")
	assert.Equal(t, routeRepo.BoundRoute.Host, "my-hostname")
	assert.Contains(t, fakeUI.Outputs[5], "OK")

	assert.Contains(t, fakeUI.Outputs[6], "Uploading my-new-app...")
	assert.Equal(t, appRepo.UploadedApp.Guid, "my-new-app-guid")
	assert.Contains(t, fakeUI.Outputs[7], "OK")
}

func TestPushingAppWhenItAlreadyExists(t *testing.T) {
	domainRepo := &testhelpers.FakeDomainRepository{}
	routeRepo := &testhelpers.FakeRouteRepository{}
	existingApp := cf.Application{Name: "existing-app", Guid: "existing-app-guid"}
	appRepo := &testhelpers.FakeApplicationRepository{AppByName: existingApp}

	fakeUI := callPush([]string{"--name", "existing-app"}, basePushConfig(), appRepo, domainRepo, routeRepo)

	assert.Contains(t, fakeUI.Outputs[0], "Uploading existing-app...")
	assert.Equal(t, appRepo.UploadedApp.Guid, "existing-app-guid")
	assert.Contains(t, fakeUI.Outputs[1], "OK")
}

func callPush(args []string,
	config *configuration.Configuration,
	appRepo api.ApplicationRepository,
	domainRepo api.DomainRepository,
	routeRepo api.RouteRepository) (fakeUI *testhelpers.FakeUI) {

	fakeUI = new(testhelpers.FakeUI)
	target := NewPush(fakeUI, config, appRepo, domainRepo, routeRepo)
	target.Run(testhelpers.NewContext("push", args))
	return
}

func basePushConfig() (config *configuration.Configuration) {
	config = testhelpers.Login()
	config.Organization = cf.Organization{Name: "MyOrg"}
	config.Space = cf.Space{Name: "MySpace"}

	return
}
