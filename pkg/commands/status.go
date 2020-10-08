package command

import (
	"fmt"
	"giterate/pkg/entities"
	"os"
	"strings"

	"github.com/gookit/color"
	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/transport/http"
	"github.com/go-git/go-git/plumbing/transport/ssh"
	"github.com/mitchellh/cli"
)

type StatusCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *StatusCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	file := initConf(arguments)
	services := entities.PopulateRepositories(file, arguments)
	Status(services, arguments)
	return 0
}

func (c *StatusCommand) Help() string {
	helpText := `
Usage: giterate status [options]

    Check status of each git repositories according to configuration file

    By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
    --config-file        set json/yaml configuration file path
    -f, --full           show status of all repositories, even if there's no uncommited changes
    --log-level          set log level (info, warn, error, debug). default: info

`
	return strings.TrimSpace(helpText)
}

func (c *StatusCommand) Synopsis() string {
	return "get the status of all repositories"
}

func Status(services []entities.Service, arguments *entities.Arguments) {
	green := color.Style{color.FgLightGreen, color.OpBold}.Render
	for _, service := range services {
		var pullOptions = git.PullOptions{RemoteName: "origin", Progress: os.Stdout}
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
			r, err := git.PlainOpen(repository.Destination)
			if err != nil {
				log.Debugf("Cannot open repository: %s, Ignoring\n", err)
			} else {
				w, err := r.Worktree()
				if err != nil {
					log.Fatalf("Cannot get worktree: %s\n", err)
				}
				head, _ := r.Head()
				currentBranch := head.Name()
				status, err := w.Status()
				if err != nil {
					log.Fatalf("Cannot get repository status: %s\n", err)
				}
				if arguments.Full || !status.IsClean() {
					fmt.Printf("\n%s: %s (%s)\n%s: %s\n%s:\n%s\n", green("Repository"), repository.URL, repository.Destination, green("Branch"), currentBranch.Short(), green("Changes"), status.String())
				}
			}
		}
	}
}
