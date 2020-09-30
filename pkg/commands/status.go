package command

import (
	"fmt"
	"giterate/pkg/entities"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/transport/http"
	"github.com/go-git/go-git/plumbing/transport/ssh"
	"github.com/mitchellh/cli"
)

type StatusCommand struct {
	Ui cli.Ui
	Arguments
}

func (c *StatusCommand) Run(args []string) int {
	file := initConf(c.Arguments.process(args))
	services := entities.PopulateRepositories(file)
	Status(services)
	return 0
}

func (c *StatusCommand) Help() string {
	helpText := `
Usage: giterate status [options]

    Check status of each git repositories according to configuration file
	
	By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
	--config-file        set json/yaml configuration file path


`
	return strings.TrimSpace(helpText)
}

func (c *StatusCommand) Synopsis() string {
	return "get the status of all repositories"
}

func Status(services []entities.Service) {
	for _, service := range services {
		var pullOptions = git.PullOptions{RemoteName: "origin", Progress: os.Stdout}
		if service.SSHPrivateKeyPath != "" {
			publicKeys, err := ssh.NewPublicKeysFromFile("git", service.SSHPrivateKeyPath, "")
			if err != nil {
				log.Fatalf("Cannot get public key: %s", err)
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
				log.Fatalf("Cannot open repository: %s", err)
			}
			w, err := r.Worktree()
			if err != nil {
				log.Fatalf("Cannot get worktree: %s", err)
			}
			fmt.Println("[INFO] Status " + repository.URL)
			status, err := w.Status()
			if err != nil {
				log.Fatalf("Cannot get repository status: %s", err)
			}
			fmt.Println(status)
		}
	}
}
