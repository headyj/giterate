package command

import (
	"fmt"
	"giterate/pkg/entities"
	"strings"

	"github.com/gookit/color"
	"github.com/mitchellh/cli"
)

type RepositoriesCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *RepositoriesCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	file := initConf(arguments)
	services := entities.PopulateRepositories(file, arguments)
	Repositories(services, arguments)
	return 0
}

func (c *RepositoriesCommand) Help() string {
	helpText := `
Usage: giterate repositories [options]

list configured repositories
	
	By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
    -p, --provider BaseURL or name    target one or multiple providers (chain multiple times)
	--config-file                     set json/yaml configuration file path
	--log-level                       set log level (info, warn, error, debug). default: info


`
	return strings.TrimSpace(helpText)
}

func (c *RepositoriesCommand) Synopsis() string {
	return "list configured repositories"
}

func Repositories(services []entities.Service, arguments *entities.Arguments) {
	green := color.Style{color.FgLightGreen, color.OpBold}.Render
	for _, service := range services {
		if service.Name != "" {
			fmt.Printf("%s (%s)\n", green(service.Name), service.BaseURL)
		} else {
			fmt.Printf("%s\n", green(service.BaseURL))
		}
		for _, repository := range service.Repositories {
			fmt.Printf("    - %s\n", green(repository.Destination))
		}
	}
}
