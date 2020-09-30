package command

import (
	"fmt"
	"giterate/pkg/entities"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/transport"
	"github.com/go-git/go-git/plumbing/transport/http"
	"github.com/go-git/go-git/plumbing/transport/ssh"
	"github.com/mitchellh/cli"
)

type CloneCommand struct {
	Ui cli.Ui
	Arguments
}

func (c *CloneCommand) Run(args []string) int {
	file := initConf(c.Arguments.process(args))
	services := entities.PopulateRepositories(file)
	Clone(services)
	//servicesJSON, _ := json.Marshal(services)
	//fmt.Printf("%s", servicesJSON)
	return 0
}

func (c *CloneCommand) Help() string {
	helpText := `
Usage: giterate clone [options]

	Clone repositories accordint to configuration file
	
	By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate/ folder

Options:
	--config-file        set json/yaml configuration file path


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
				log.Fatalf("Cannot get public key: %s", err)
			}
			cloneOptions.Auth = publicKeys
		} else if service.Username != "" && service.Password != "" {
			cloneOptions.Auth = &http.BasicAuth{
				Username: service.Username,
				Password: service.Password,
			}
		}
		for _, repository := range service.Repositories {
			cloneOptions.URL = repository.URL
			if repository.CloneOptions != nil {
				entities.ProcessCloneOptions(&repository.CloneOptions, &cloneOptions)
			}
			fmt.Println("[INFO] Cloning " + repository.URL + "...")
			_, err := git.PlainClone(repository.Destination, false, &cloneOptions)
			if err != nil {
				if err == transport.ErrEmptyRemoteRepository || err == git.ErrRepositoryAlreadyExists {
					fmt.Println("[INFO] Repository already cloned, ignoring...")
					continue
				} else {
					log.Fatalf("Cannot clone repository: %s", err)
				}
			} else {
				fmt.Println("[INFO] Cloned " + repository.URL + "...")
			}
		}
	}
}
