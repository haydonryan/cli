package application

import (
	"cf"
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"strings"
	"time"
	"github.com/cloudfoundry/loggregatorlib/logmessage"
)

type Start struct {
	ui        terminal.UI
	config    *configuration.Configuration
	appRepo   api.ApplicationRepository
	logRepo   api.LogsRepository
	startTime time.Time
	appReq    requirements.ApplicationRequirement
}

type ApplicationStarter interface {
	ApplicationStart(cf.Application) (startedApp cf.Application, err error)
}

func NewStart(ui terminal.UI, config *configuration.Configuration, appRepo api.ApplicationRepository, logRepo api.LogsRepository) (cmd *Start) {
	cmd = new(Start)
	cmd.ui = ui
	cmd.config = config
	cmd.appRepo = appRepo
	cmd.logRepo = logRepo

	return
}

func (cmd *Start) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) == 0 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "start")
		return
	}

	cmd.appReq = reqFactory.NewApplicationRequirement(c.Args()[0])

	reqs = []requirements.Requirement{cmd.appReq}
	return
}

func (cmd *Start) Run(c *cli.Context) {
	cmd.ApplicationStart(cmd.appReq.GetApplication())
}

func (cmd *Start) ApplicationStart(app cf.Application) (updatedApp cf.Application, err error) {
	if app.State == "started" {
		cmd.ui.Say(terminal.WarningColor("App " + app.Name + " is already started"))
		return
	}

	cmd.ui.Say("Starting app %s in org %s / space %s as %s...",
		terminal.EntityNameColor(app.Name),
		terminal.EntityNameColor(cmd.config.Organization.Name),
		terminal.EntityNameColor(cmd.config.Space.Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	updatedApp, apiResponse := cmd.appRepo.Start(app)
	if apiResponse.IsNotSuccessful() {
		cmd.ui.Failed(apiResponse.Message)
		return
	}

	cmd.ui.Ok()

	onConnect := func() {}
	onMessage := func(msg *logmessage.Message) {
		cmd.ui.Say(logMessageOutput(msg))
	}

	go cmd.logRepo.TailLogsFor(app, onConnect, onMessage, 1)
//	cmd.logRepo.TailLogsFor(app, onConnect, onMessage, 0)

	instances, apiResponse := cmd.appRepo.GetInstances(app)

	for apiResponse.IsNotSuccessful() {
		if apiResponse.ErrorCode != cf.APP_NOT_STAGED {
			cmd.ui.Say("")
			cmd.ui.Failed(apiResponse.Message)
			return
		}

		cmd.ui.Wait(1 * time.Second)
		instances, apiResponse = cmd.appRepo.GetInstances(app)
		cmd.ui.LoadingIndication()
	}

	cmd.ui.Say("")

	cmd.startTime = time.Now()

	for cmd.displayInstancesStatus(app, instances) {
		cmd.ui.Wait(1 * time.Second)
		instances, _ = cmd.appRepo.GetInstances(app)
	}

	return
}

func (cmd Start) displayInstancesStatus(app cf.Application, instances []cf.ApplicationInstance) (notFinished bool) {
	totalCount := len(instances)
	runningCount, startingCount, flappingCount, downCount := 0, 0, 0, 0

	for _, inst := range instances {
		switch inst.State {
		case cf.InstanceRunning:
			runningCount++
		case cf.InstanceStarting:
			startingCount++
		case cf.InstanceFlapping:
			flappingCount++
		case cf.InstanceDown:
			downCount++
		}
	}

	if flappingCount > 0 {
		cmd.ui.Failed("Start unsuccessful")
		return false
	}

	anyInstanceRunning := runningCount > 0

	if anyInstanceRunning {
		if len(app.Urls) == 0 {
			cmd.ui.Say("Started")
		} else {
			cmd.ui.Say("Started: app %s available at %s", app.Name, app.Urls[0])
		}
		return false
	} else {
		details := instancesDetails(runningCount, startingCount, downCount)
		cmd.ui.Say("%d of %d instances running (%s)", runningCount, totalCount, details)
	}

	if time.Since(cmd.startTime) > cmd.config.ApplicationStartTimeout*time.Second {
		cmd.ui.Failed("Start app timeout")
		return false
	}

	return totalCount > runningCount
}

func instancesDetails(runningCount int, startingCount int, downCount int) string {
	details := []string{}

	if startingCount > 0 {
		details = append(details, fmt.Sprintf("%d starting", startingCount))
	}

	if downCount > 0 {
		details = append(details, fmt.Sprintf("%d down", downCount))
	}

	return strings.Join(details, ", ")
}
