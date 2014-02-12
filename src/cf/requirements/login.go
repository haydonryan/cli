package requirements

import (
	"cf/configuration"
	"cf/terminal"
)

type LoginRequirement struct {
	ui     terminal.UI
	config configuration.Reader
}

func NewLoginRequirement(ui terminal.UI, config configuration.Reader) LoginRequirement {
	return LoginRequirement{ui, config}
}

func (req LoginRequirement) Execute() (success bool) {
	if !req.config.IsLoggedIn() {
		req.ui.Say(terminal.NotLoggedInText())
		return false
	}
	return true
}
