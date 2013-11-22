package application_test

import (
	"cf"
	. "cf/commands/application"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
	"time"
)

func TestEventsRequirements(t *testing.T) {
	reqFactory, eventsRepo := getEventsDependencies()

	callEvents(t, []string{"my-app"}, reqFactory, eventsRepo)

	assert.Equal(t, reqFactory.ApplicationName, "my-app")
	assert.True(t, testcmd.CommandDidPassRequirements)
}

func TestEventsFailsWithUsage(t *testing.T) {
	reqFactory, eventsRepo := getEventsDependencies()
	ui := callEvents(t, []string{}, reqFactory, eventsRepo)

	assert.True(t, ui.FailedWithUsage)
	assert.False(t, testcmd.CommandDidPassRequirements)
}

func TestEventsSuccess(t *testing.T) {
	timestamp, err := time.Parse(TIMESTAMP_FORMAT, "2000-01-01T00:01:11.00-0000")
	assert.NoError(t, err)

	reqFactory, eventsRepo := getEventsDependencies()
	app := cf.Application{}
	app.Name = "my-app"
	reqFactory.Application = app

	eventsRepo.Events = []cf.EventFields{
		{
			InstanceIndex:   98,
			Timestamp:       timestamp,
			ExitDescription: "app instance exited",
			ExitStatus:      78,
		},
		{
			InstanceIndex:   99,
			Timestamp:       timestamp,
			ExitDescription: "app instance was stopped",
			ExitStatus:      77,
		},
	}

	ui := callEvents(t, []string{"my-app"}, reqFactory, eventsRepo)

	assert.Contains(t, ui.Outputs[0], "Getting events for app")
	assert.Contains(t, ui.Outputs[0], "my-app")
	assert.Contains(t, ui.Outputs[0], "my-org")
	assert.Contains(t, ui.Outputs[0], "my-space")
	assert.Contains(t, ui.Outputs[0], "my-user")
	assert.Contains(t, ui.Outputs[1], "time")
	assert.Contains(t, ui.Outputs[1], "instance")
	assert.Contains(t, ui.Outputs[1], "description")
	assert.Contains(t, ui.Outputs[1], "exit status")
	assert.Contains(t, ui.Outputs[2], timestamp.Local().Format(TIMESTAMP_FORMAT))
	assert.Contains(t, ui.Outputs[2], "98")
	assert.Contains(t, ui.Outputs[2], "app instance exited")
	assert.Contains(t, ui.Outputs[2], "78")
	assert.Contains(t, ui.Outputs[3], timestamp.Local().Format(TIMESTAMP_FORMAT))
	assert.Contains(t, ui.Outputs[3], "99")
	assert.Contains(t, ui.Outputs[3], "app instance was stopped")
	assert.Contains(t, ui.Outputs[3], "77")
}

func TestEventsWhenNoEventsAvailable(t *testing.T) {
	reqFactory, eventsRepo := getEventsDependencies()
	app := cf.Application{}
	app.Name = "my-app"
	reqFactory.Application = app

	ui := callEvents(t, []string{"my-app"}, reqFactory, eventsRepo)

	assert.Contains(t, ui.Outputs[0], "events")
	assert.Contains(t, ui.Outputs[0], "my-app")
	assert.Contains(t, ui.Outputs[1], "No events")
	assert.Contains(t, ui.Outputs[1], "my-app")
}

func getEventsDependencies() (reqFactory *testreq.FakeReqFactory, eventsRepo *testapi.FakeAppEventsRepo) {
	reqFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true}
	eventsRepo = &testapi.FakeAppEventsRepo{}
	return
}

func callEvents(t *testing.T, args []string, reqFactory *testreq.FakeReqFactory, eventsRepo *testapi.FakeAppEventsRepo) (ui *testterm.FakeUI) {
	ui = new(testterm.FakeUI)
	ctxt := testcmd.NewContext("events", args)

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

	cmd := NewEvents(ui, config, eventsRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
