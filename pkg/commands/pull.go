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

type PullCommand struct {
	Ui cli.Ui
	Arguments
}

func (c *PullCommand) Run(args []string) int {
	file := initConf(c.Arguments.process(args))
	services := entities.PopulateRepositories(file)
	Pull(services)
	return 0
}

func (c *PullCommand) Help() string {
	helpText := `
Usage: giterate pull [options]

    Pull repositories on current branches according to configuration file
	
	By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
	--config-file        set json/yaml configuration file path


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
			if repository.CloneOptions != nil {
				entities.ProcessPullOptions(&repository.CloneOptions, &pullOptions)
			}
			r, err := git.PlainOpen(repository.Destination)
			if err != nil {
				log.Fatalf("Cannot open repository: %s", err)
			}
			w, err := r.Worktree()
			if err != nil {
				log.Fatalf("Cannot get worktree: %s", err)
			}
			fmt.Println("[INFO] Pull " + repository.URL)
			err = w.Pull(&pullOptions)
			switch err {
			case git.NoErrAlreadyUpToDate:
				continue
			case git.ErrUnstagedChanges:
				fmt.Println("[INFO] Repository contains unstaged changes, ignoring")
			default:
				log.Fatalf("Cannot pull respository: %s", err)
			}
		}
	}
}
