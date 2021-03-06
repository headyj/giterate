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

type PullCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *PullCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	file := initConf(arguments)
	services := entities.PopulateRepositories(file, arguments)
	Pull(services)
	return 0
}

func (c *PullCommand) Help() string {
	helpText := `
Usage: giterate pull [options]

    Pull repositories on current branches according to configuration file

    By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
    -r, --repository path             target one or multiple repositories (chain multiple times)
    -p, --provider BaseURL or name    target one or multiple providers (chain multiple times)
    --config-file                     set json/yaml configuration file path
    --log-level                       set log level (info, warn, error, debug). default: info

`
	return strings.TrimSpace(helpText)
}

func (c *PullCommand) Synopsis() string {
	return "pull repositories on current branches"
}

func Pull(services []entities.Service) {
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
			if repository.CloneOptions != nil {
				entities.ProcessPullOptions(&repository.CloneOptions, &pullOptions)
			}
			r, err := git.PlainOpen(repository.Destination)
			if err != nil {
				log.Errorf("Cannot open repository: %s, ignoring\n", err)
			} else {
				w, err := r.Worktree()
				if err != nil {
					log.Errorf("Cannot get worktree: %s, ignoring\n", err)
				} else {
					log.Infof("Pull %s\n", repository.URL)
					err = w.Pull(&pullOptions)
					switch err {
					case git.NoErrAlreadyUpToDate:
						continue
					case git.ErrUnstagedChanges:
						log.Info("Repository contains unstaged changes, ignoring\n")
					default:
						log.Errorf("Cannot pull respository: %s, ignoring\n", err)
					}
				}
			}
		}
	}
}
