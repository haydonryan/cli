package main

import (
	"os"
	"cf/app"
	"cf/requirements"
	"cf/commands"
	"cf/api"
	"cf/terminal"
	"cf/configuration"
	"github.com/codegangsta/cli"
)

func main() {
	termUI := new(terminal.TerminalUI)
	assignTemplates()
	config := loadConfig(termUI)
	repoLocator := api.NewRepositoryLocator(config)
	cmdFactory := commands.NewFactory(termUI, repoLocator)
	reqFactory := requirements.NewFactory(termUI, repoLocator)

	app, err := app.NewApp(cmdFactory, reqFactory)
	if err != nil {
		return
	}
	app.Run(os.Args)
}

func assignTemplates() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   [environment variables] {{.Name}} [global options] command [arguments...] [command options]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Description}}
   {{end}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
ENVIRONMENT VARIABLES:
   CF_TRACE=true - will output HTTP requests and responses during command
   HTTP_PROXY=http://proxy.example.com:8080 - set to your proxy
`

	cli.CommandHelpTemplate = `NAME:
   {{.Name}} - {{.Description}}
{{with .ShortName}}
ALIAS:
   {{.}}
{{end}}
USAGE:
   {{.Usage}}{{with .Flags}}

OPTIONS:
   {{range .}}{{.}}
   {{end}}{{else}}
{{end}}`

}

func loadConfig(termUI terminal.UI) (config configuration.Configuration) {
	configRepo := configuration.NewConfigurationDiskRepository()
	config, err := configRepo.Get()
	if err != nil {
		termUI.ConfigFailure(err)
		return
	}
	return
}
