package service

import (
	"cf"
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/codegangsta/cli"
)

const AppAlreadyBoundErrorCode = "90003"

type BindService struct {
	ui                 terminal.UI
	config             *configuration.Configuration
	serviceBindingRepo api.ServiceBindingRepository
	appReq             requirements.ApplicationRequirement
	serviceInstanceReq requirements.ServiceInstanceRequirement
}

func NewBindService(ui terminal.UI, config *configuration.Configuration, serviceBindingRepo api.ServiceBindingRepository) (cmd *BindService) {
	cmd = new(BindService)
	cmd.ui = ui
	cmd.config = config
	cmd.serviceBindingRepo = serviceBindingRepo
	return
}

func (cmd *BindService) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {

	if len(c.Args()) != 2 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "bind-service")
		return
	}
	appName := c.Args()[0]
	serviceName := c.Args()[1]

	cmd.appReq = reqFactory.NewApplicationRequirement(appName)
	cmd.serviceInstanceReq = reqFactory.NewServiceInstanceRequirement(serviceName)

	reqs = []requirements.Requirement{cmd.appReq, cmd.serviceInstanceReq}
	return
}

func (cmd *BindService) Run(c *cli.Context) {
	app := cmd.appReq.GetApplication()
	instance := cmd.serviceInstanceReq.GetServiceInstance()

	cmd.BindApplication(app, instance)
}

func (cmd *BindService) BindApplication(app cf.Application, serviceInstance cf.ServiceInstance) {
	cmd.ui.Say("Binding service %s to app %s in org %s / space %s as %s...",
		terminal.EntityNameColor(serviceInstance.Name),
		terminal.EntityNameColor(app.Name),
		terminal.EntityNameColor(cmd.config.OrganizationFields.Name),
		terminal.EntityNameColor(cmd.config.SpaceFields.Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	apiResponse := cmd.serviceBindingRepo.Create(serviceInstance.Guid, app.Guid)
	if apiResponse.IsNotSuccessful() && apiResponse.ErrorCode != "90003" {
		cmd.ui.Failed(apiResponse.Message)
		return
	}

	cmd.ui.Ok()

	if apiResponse.ErrorCode == AppAlreadyBoundErrorCode {
		cmd.ui.Warn("App %s is already bound to %s.", app.Name, serviceInstance.Name)
		return
	}

	cmd.ui.Say("TIP: Use '%s push' to ensure your env variable changes take effect", cf.Name())
}
