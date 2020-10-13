package command

import (
	"fmt"
	"giterate/pkg/entities"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/transport/http"
	"github.com/go-git/go-git/plumbing/transport/ssh"
	"github.com/mitchellh/cli"
)

type PushCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *PushCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	file := initConf(c.Arguments.Process(args))
	services := entities.PopulateRepositories(file, arguments)
	Push(services)
	return 0
}

func (c *PushCommand) Help() string {
	helpText := `
Usage: giterate push [options]

    Push commited changes

    By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate/ folder

Options:
    -r, --repository path             target one or multiple repositories (chain multiple times)
    -p, --provider BaseURL or name    target one or multiple providers (chain multiple times)
    --config-file                     set json/yaml configuration file path
    --log-level                       set log level (info, warn, error, debug). default: info

`

	return strings.TrimSpace(helpText)
}

func (c *PushCommand) Synopsis() string {
	return "push commited changes"
}

func Push(services []entities.Service) {
	for _, service := range services {
		var pushOptions = git.PushOptions{Progress: os.Stdout}
		if service.SSHPrivateKeyPath != "" {
			fmt.Println("ssh login")
			publicKeys, err := ssh.NewPublicKeysFromFile("git", service.SSHPrivateKeyPath, "")
			if err != nil {
				log.Fatalf("Cannot get public key: %s\n", err)
			}
			pushOptions.Auth = publicKeys
		} else if service.Username != "" && service.Password != "" {
			pushOptions.Auth = &http.BasicAuth{
				Username: service.Username,
				Password: service.Password,
			}
		}
		for _, repository := range service.Repositories {
			r, err := git.PlainOpen(repository.Destination)
			if repository.CloneOptions != nil {
				entities.ProcessPushOptions(&repository.CloneOptions, &pushOptions)
			}
			err = r.Push(&pushOptions)
			if err != nil {
				switch err {
				case git.NoErrAlreadyUpToDate:
					continue
				default:
					log.Errorf("Cannot push to repository: %s\n", err)
				}
			}

			log.Infof("Pushed %s\n", repository.URL)
		}
	}
}
