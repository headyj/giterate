package command

import (
	"fmt"
	"giterate/pkg/entities"
	"strings"

	"github.com/gookit/color"
	"github.com/mitchellh/cli"
)

type ProvidersCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *ProvidersCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	services := initConf(arguments)
	Providers(services, arguments)
	return 0
}

func (c *ProvidersCommand) Help() string {
	helpText := `
Usage: giterate providers [options]

list configured providers
	
	By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
	--config-file        set json/yaml configuration file path
	--log-level          set log level (info, warn, error, debug). default: info


`
	return strings.TrimSpace(helpText)
}

func (c *ProvidersCommand) Synopsis() string {
	return "list configured providers"
}

func Providers(services []entities.Service, arguments *entities.Arguments) {
	green := color.Style{color.FgLightGreen, color.OpBold}.Render
	for _, service := range services {
		if service.Name != "" {
			fmt.Printf("%s (%s)\n", green(service.Name), service.BaseURL)
		} else {
			fmt.Printf("%s\n", green(service.BaseURL))
		}
	}
}
