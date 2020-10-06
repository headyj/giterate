package command

import (
	"fmt"
	"giterate/pkg/entities"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing"
	"github.com/mitchellh/cli"
)

type CheckoutCommand struct {
	Ui cli.Ui
	entities.Arguments
}

func (c *CheckoutCommand) Run(args []string) int {
	arguments := c.Arguments.Process(args)
	file := initConf(arguments)
	services := entities.PopulateRepositories(file, arguments)
	Checkout(services, arguments)
	//servicesJSON, _ := json.Marshal(services)
	//fmt.Printf("%s", servicesJSON)
	return 0
}

func (c *CheckoutCommand) Help() string {
	helpText := `
Usage: giterate checkout [options]

checkout the configured branch on all repositories
	
	By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder

Options:
    -r, --repository path             target one or multiple repositories (chain multiple times)
    -p, --provider BaseURL or name    target one or multiple providers (chain multiple times)
    --force                           reset uncommited changes
    --config-file                     set json/yaml configuration file path
    --log-level                       set log level (info, warn, error, debug). default: info


`
	return strings.TrimSpace(helpText)
}

func (c *CheckoutCommand) Synopsis() string {
	return "checkout the configured branch on all repositories"
}

func Checkout(services []entities.Service, arguments *entities.Arguments) {
	for _, service := range services {
		var checkoutOptions = git.CheckoutOptions{}
		if arguments.Force {
			checkoutOptions.Force = true
		}
		for _, repository := range service.Repositories {
			log.Debugf("Checkout branch %s for %s\n", checkoutOptions.Branch.String(), repository.URL)
			checkoutOptions.Branch = plumbing.NewBranchReferenceName(repository.DefaultBranch)
			if repository.CloneOptions != nil {
				entities.ProcessCheckoutOptions(&repository.CloneOptions, &checkoutOptions)
			}
			r, err := git.PlainOpen(repository.Destination)
			if err != nil {
				log.Infof("Cannot open repository: %s, ignoring\n", err)
			} else {
				w, err := r.Worktree()
				if err != nil {
					log.Fatalf("Cannot get worktree: %s\n", err)
				}
				log.Infof("Checkout branch %s for %s\n", checkoutOptions.Branch.String(), repository.URL)
				err = w.Checkout(&checkoutOptions)
				if err != nil {
					switch err {
					case git.ErrUnstagedChanges:
						log.Info("Worktree contains unstaged changes, ignoring\n")
					default:
						log.Infof("Cannot checkout branch: %s. Ignoring\n", err)
						headRef, _ := r.Head()
						fmt.Println(headRef)
					}
				}
			}
		}
	}
}
