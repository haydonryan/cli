package commands

import (
	"cf"
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"fmt"
	"github.com/codegangsta/cli"
)

type InfoResponse struct {
	ApiVersion            string `json:"api_version"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

type Target struct {
	ui         terminal.UI
	config     *configuration.Configuration
	configRepo configuration.ConfigurationRepository
	orgRepo    api.OrganizationRepository
	spaceRepo  api.SpaceRepository
}

func NewTarget(ui terminal.UI, configRepo configuration.ConfigurationRepository, orgRepo api.OrganizationRepository, spaceRepo api.SpaceRepository) (t Target) {
	t.ui = ui
	t.configRepo = configRepo
	t.config, _ = configRepo.Get()
	t.orgRepo = orgRepo
	t.spaceRepo = spaceRepo

	return
}

func (cmd Target) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
	}
	return
}

func (t Target) Run(c *cli.Context) {
	argsCount := len(c.Args())
	orgName := c.String("o")
	spaceName := c.String("s")

	if argsCount == 0 && orgName == "" && spaceName == "" {
		t.ui.ShowConfiguration(t.config)

		if !t.config.HasOrganization() {
			t.ui.Say("No org targeted. Use '%s target -o' to target an org.", cf.Name)
		}
		if !t.config.HasSpace() {
			t.ui.Say("No space targeted. Use '%s target -s' to target a space.", cf.Name)
		}
		return
	}

	if orgName != "" {
		t.setOrganization(orgName)
		if t.config.IsLoggedIn() {
			t.ui.Say("No space targeted. Use '%s target -s' to target a space.", cf.Name)
		}
		return
	}

	if spaceName != "" {
		t.setSpace(spaceName)
		return
	}

	return
}

func (t Target) setOrganization(orgName string) {
	if !t.config.IsLoggedIn() {
		t.ui.Failed("You must be logged in to set an organization. Use '%s login'.", cf.Name)
		return
	}

	org, found, err := t.orgRepo.FindByName(orgName)
	if err != nil {
		t.ui.Failed("Could not set organization.")
		return
	}

	if !found {
		t.ui.Failed(fmt.Sprintf("Organization %s not found.", orgName))
		return
	}

	t.config.Organization = org
	t.config.Space = cf.Space{}
	t.saveAndShowConfig()
}

func (t Target) setSpace(spaceName string) {
	if !t.config.IsLoggedIn() {
		t.ui.Failed("You must be logged in to set a space. Use '%s login'.", cf.Name)
		return
	}

	if !t.config.HasOrganization() {
		t.ui.Failed("Organization must be set before targeting space.")
		return
	}

	space, found, err := t.spaceRepo.FindByName(spaceName)
	if err != nil {
		t.ui.Failed("You do not have access to that space.")
		return
	}

	if !found {
		t.ui.Failed(fmt.Sprintf("Space %s not found.", spaceName))
		return
	}

	t.config.Space = space
	t.saveAndShowConfig()
}

func (t Target) saveAndShowConfig() {
	err := t.configRepo.Save()
	if err != nil {
		t.ui.Failed(err.Error())
		return
	}
	t.ui.ShowConfiguration(t.config)
}
