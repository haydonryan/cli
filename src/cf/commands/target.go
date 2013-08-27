package commands

import (
	"cf/api"
	"cf/configuration"
	term "cf/terminal"
	"github.com/codegangsta/cli"
)

type InfoResponse struct {
	ApiVersion            string `json:"api_version"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

type Target struct {
	ui        term.UI
	config    *configuration.Configuration
	orgRepo   api.OrganizationRepository
	spaceRepo api.SpaceRepository
}

func NewTarget(ui term.UI, config *configuration.Configuration, orgRepo api.OrganizationRepository, spaceRepo api.SpaceRepository) (t Target) {
	t.ui = ui
	t.config = config
	t.orgRepo = orgRepo
	t.spaceRepo = spaceRepo

	return
}

func (t Target) Run(c *cli.Context) {
	argsCount := len(c.Args())
	orgName := c.String("o")
	spaceName := c.String("s")

	if argsCount == 0 && orgName == "" && spaceName == "" {
		t.ui.ShowConfiguration(t.config)
		return
	}

	if argsCount > 0 {
		t.setNewTarget(c.Args()[0])
		return
	}

	if orgName != "" {
		t.setOrganization(orgName)
		return
	}

	if spaceName != "" {
		t.setSpace(spaceName)
		return
	}

	return
}

func (t Target) setNewTarget(target string) {
	url := "https://" + target
	t.ui.Say("Setting target to %s...", term.Yellow(url))

	request, err := api.NewAuthorizedRequest("GET", url+"/v2/info", "", nil)

	if err != nil {
		t.ui.Failed("URL invalid.", err)
		return
	}

	serverResponse := new(InfoResponse)
	err = api.PerformRequestAndParseResponse(request, &serverResponse)

	if err != nil {
		t.ui.Failed("", err)
		return
	}

	err = t.saveTarget(url, serverResponse)

	if err != nil {
		t.ui.Failed("Error saving configuration", err)
		return
	}

	t.ui.Ok()
	t.ui.ShowConfiguration(t.config)
}

func (t *Target) saveTarget(target string, info *InfoResponse) (err error) {
	t.config = new(configuration.Configuration)
	t.config.Target = target
	t.config.ApiVersion = info.ApiVersion
	t.config.AuthorizationEndpoint = info.AuthorizationEndpoint
	err = t.config.Save()
	return
}

func (t Target) setOrganization(orgName string) {
	if !t.config.IsLoggedIn() {
		t.ui.Failed("You must be logged in to set an organization.", nil)
		return
	}

	org, err := t.orgRepo.FindByName(t.config, orgName)
	if err != nil {
		t.ui.Failed("Could not set organization.", nil)
		return
	}

	t.config.Organization = org
	t.saveAndShowConfig()
}

func (t Target) setSpace(spaceName string) {
	if !t.config.IsLoggedIn() {
		t.ui.Failed("You must be logged in to set a space.", nil)
		return
	}

	if !t.config.HasOrganization() {
		t.ui.Failed("Organization must be set before targeting space.", nil)
		return
	}

	space, err := t.spaceRepo.FindByName(t.config, spaceName)
	if err != nil {
		t.ui.Failed("You do not have access to that space.", nil)
		return
	}

	t.config.Space = space
	t.saveAndShowConfig()
}

func (t Target) saveAndShowConfig() {
	err := t.config.Save()
	if err != nil {
		t.ui.Failed("Error saving configuration", err)
		return
	}
	t.ui.ShowConfiguration(t.config)
}
