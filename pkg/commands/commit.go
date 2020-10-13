package command

import (
	"bufio"
	"fmt"
	"giterate/pkg/entities"
	"os"
	"strings"

	"github.com/go-git/go-git"
	"github.com/gookit/color"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type CommitCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *CommitCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	file := initConf(arguments)
	services := entities.PopulateRepositories(file, arguments)
	Commit(services, arguments)
	return 0
}

func (c *CommitCommand) Help() string {
	helpText := `
Usage: giterate commit [options]

    Check changes and ask for commit message in case of changes
	
    By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
    -r, --repository path             target one or multiple repositories (chain multiple times)
    -p, --provider BaseURL or name    target one or multiple providers (chain multiple times)
    --config-file                     set json/yaml configuration file path
    --log-level                       set log level (info, warn, error, debug). default: info


`
	return strings.TrimSpace(helpText)
}

func (c *CommitCommand) Synopsis() string {
	return "check changes and ask for commit message in case of changes"
}

func Commit(services []entities.Service, arguments *entities.Arguments) {
	green := color.Style{color.FgLightGreen, color.OpBold}.Render
	for _, service := range services {
		for _, repository := range service.Repositories {
			r, err := git.PlainOpen(repository.Destination)
			if err != nil {
				log.Debugf("Cannot open repository: %s, ignoring\n", err)
			} else {
				w, err := r.Worktree()
				if err != nil {
					log.Fatalf("Cannot get worktree: %s\n", err)
				}
				head, err := r.Head()
				if err != nil {
					log.Fatalf("Cannot get repository HEAD reference: %s\n", err)
				}
				currentBranch := head.Name()
				status, err := w.Status()
				if err != nil {
					log.Fatalf("Cannot get status: %s\n", err)
				}
				if !status.IsClean() {
					reader := bufio.NewReader(os.Stdin)
					fmt.Printf("\n\n%s: %s (%s)\n%s: %s\n%s:\n%s\nEnter commit message (let empty to ignore):", green("Repository"), repository.URL, repository.Destination, green("Branch"), currentBranch.Short(), green("Changes"), status.String())
					commitMsg, err := reader.ReadString('\n')
					if err != nil {
						log.Fatalf("Cannot read commit message: %s\n", err)
					}
					if commitMsg != "\n" {
						_, err = w.Add(".")
						commit, err := w.Commit(commitMsg, &git.CommitOptions{})
						if err != nil {
							log.Fatalf("Cannot commit to repository: %s\n", err)
						}
						_, err = r.CommitObject(commit)
						if err != nil {
							log.Fatalf("Cannot commit to repository: %s\n", err)
						}
					}
				}
			}
		}
	}
}
