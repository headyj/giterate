package command

import (
	"giterate/pkg/entities"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/transport"
	"github.com/go-git/go-git/plumbing/transport/http"
	"github.com/go-git/go-git/plumbing/transport/ssh"
	"github.com/mitchellh/cli"
)

type CloneCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *CloneCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	file := initConf(c.Arguments.Process(args))
	services := entities.PopulateRepositories(file, arguments)
	Clone(services)
	//servicesJSON, _ := json.Marshal(services)
	//fmt.Printf("%s", servicesJSON)
	return 0
}

func (c *CloneCommand) Help() string {
	helpText := `
Usage: giterate clone [options]

    Clone repositories according to configuration file

    By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate/ folder

Options:
    -r, --repository path             target one or multiple repositories (chain multiple times)
    -p, --provider BaseURL or name    target one or multiple providers (chain multiple times)
    --config-file                     set json/yaml configuration file path
    --log-level                       set log level (info, warn, error, debug). default: info

`

	return strings.TrimSpace(helpText)
}

func (c *CloneCommand) Synopsis() string {
	return "clone repositories according to configuration file"
}

func Clone(services []entities.Service) {
	for _, service := range services {
		var cloneOptions = git.CloneOptions{Progress: os.Stdout}
		if service.SSHPrivateKeyPath != "" {
			publicKeys, err := ssh.NewPublicKeysFromFile("git", service.SSHPrivateKeyPath, "")
			if err != nil {
				log.Fatalf("Cannot get public key: %s\n", err)
			}
			cloneOptions.Auth = publicKeys
		} else if service.Username != "" && service.Password != "" {
			cloneOptions.Auth = &http.BasicAuth{
				Username: service.Username,
				Password: service.Password,
			}
		}
		for _, repository := range service.Repositories {
			rCloneOptions := cloneOptions
			rCloneOptions.URL = repository.URL
			if repository.CloneOptions != nil {
				entities.ProcessCloneOptions(&repository.CloneOptions, &rCloneOptions)
			}
			log.Infof("Cloning %s...\n", repository.URL)
			_, err := git.PlainClone(repository.Destination, false, &rCloneOptions)
			if err != nil {
				if err == transport.ErrEmptyRemoteRepository || err == git.ErrRepositoryAlreadyExists {
					log.Info("Repository already cloned, ignoring...\n")
				} else {
					log.Fatalf("Cannot clone repository: %s\n", err)
				}
			} else {
				log.Infof("Cloned %s\n", repository.URL)
			}
		}
	}
}
