package command

import (
	"giterate/pkg/entities"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/transport/http"
	"github.com/go-git/go-git/plumbing/transport/ssh"
	"github.com/mitchellh/cli"
)

type ExecCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *ExecCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	file := initConf(arguments)
	services := entities.PopulateRepositories(file, arguments)
	Exec(services, arguments)
	return 0
}

func (c *ExecCommand) Help() string {
	helpText := `
Usage: giterate exec [options]

    Execute a custom git command

    By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
    -c, --command                     quoted command to be executed
    -r, --repository path             target one or multiple repositories (chain multiple times)
    -p, --provider BaseURL or name    target one or multiple providers (chain multiple times)
    --config-file                     set json/yaml configuration file path
    --log-level                       set log level (info, warn, error, debug). default: info

`
	return strings.TrimSpace(helpText)
}

func (c *ExecCommand) Synopsis() string {
	return "pull repositories on current branches"
}

func Exec(services []entities.Service, arguments *entities.Arguments) {
	for _, service := range services {
		var pullOptions = git.PullOptions{Progress: os.Stdout}
		if service.SSHPrivateKeyPath != "" {
			publicKeys, err := ssh.NewPublicKeysFromFile("git", service.SSHPrivateKeyPath, "")
			if err != nil {
				log.Fatalf("Cannot get public key: %s\n", err)
			}
			pullOptions.Auth = publicKeys
		} else if service.Username != "" && service.Password != "" {
			pullOptions.Auth = &http.BasicAuth{
				Username: service.Username,
				Password: service.Password,
			}
		}
		for _, repository := range service.Repositories {
			if len(arguments.Command) == 0 {
				log.Fatalf("Command is empty")
			}
			log.Infof("Execute command %s on %s\n", arguments.Command, repository.URL)
			cmd := exec.Command("git", arguments.Command...)
			cmd.Dir = repository.Destination
			err := cmd.Run()
			if err != nil {
				log.Errorf("Cannot execute command: %s, ignoring\n", err)
			}
		}
	}
}
