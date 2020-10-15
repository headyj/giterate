package command

import (
	"giterate/pkg/entities"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/transport/http"
	"github.com/go-git/go-git/plumbing/transport/ssh"
	"github.com/mitchellh/cli"
)

type FetchCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *FetchCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	file := initConf(arguments)
	services := entities.PopulateRepositories(file, arguments)
	Fetch(services)
	return 0
}

func (c *FetchCommand) Help() string {
	helpText := `
Usage: giterate fetch [options]

    Fetch repositories on current branches according to configuration file

    By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
    -r, --repository path             target one or multiple repositories (chain multiple times)
    -p, --provider BaseURL or name    target one or multiple providers (chain multiple times)
    --config-file                     set json/yaml configuration file path
    --log-level                       set log level (info, warn, error, debug). default: info

`
	return strings.TrimSpace(helpText)
}

func (c *FetchCommand) Synopsis() string {
	return "fetch repositories on current branches"
}

func Fetch(services []entities.Service) {
	for _, service := range services {
		var fetchOptions = git.FetchOptions{Progress: os.Stdout}
		if service.SSHPrivateKeyPath != "" {
			publicKeys, err := ssh.NewPublicKeysFromFile("git", service.SSHPrivateKeyPath, "")
			if err != nil {
				log.Fatalf("Cannot get public key: %s\n", err)
			}
			fetchOptions.Auth = publicKeys
		} else if service.Username != "" && service.Password != "" {
			fetchOptions.Auth = &http.BasicAuth{
				Username: service.Username,
				Password: service.Password,
			}
		}
		for _, repository := range service.Repositories {
			if repository.CloneOptions != nil {
				entities.ProcessFetchOptions(&repository.CloneOptions, &fetchOptions)
			}
			r, err := git.PlainOpen(repository.Destination)
			if err != nil {
				log.Errorf("Cannot open repository: %s, ignoring\n", err)
			} else {
				log.Infof("Fetch %s\n", repository.URL)
				err = r.Fetch(&fetchOptions)
				switch err {
				case git.NoErrAlreadyUpToDate:
					continue
				case git.ErrUnstagedChanges:
					log.Info("Repository contains unstaged changes, ignoring\n")
				default:
					log.Errorf("Cannot fetch respository: %s, ignoring\n", err)
				}
			}
		}
	}
}
