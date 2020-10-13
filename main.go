package main

import (
	command "giterate/pkg/commands"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	args := os.Args[1:]
	cli := cli.CLI{
		Args:     args,
		Commands: Commands,
		HelpFunc: cli.BasicHelpFunc("giterate"),
		Version:  "0.4.2",
	}

	status, err := cli.Run()
	if err != nil {
		log.Fatalf("Cannot run command: %s\n", err)
	}

	os.Exit(status)

}

var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{
		"clone": func() (cli.Command, error) {
			return &command.CloneCommand{
				Ui: ui,
			}, nil
		},
		"pull": func() (cli.Command, error) {
			return &command.PullCommand{
				Ui: ui,
			}, nil
		},
		"status": func() (cli.Command, error) {
			return &command.StatusCommand{
				Ui: ui,
			}, nil
		},
		"checkout": func() (cli.Command, error) {
			return &command.CheckoutCommand{
				Ui: ui,
			}, nil
		},
		"commit": func() (cli.Command, error) {
			return &command.CommitCommand{
				Ui: ui,
			}, nil
		},
		"push": func() (cli.Command, error) {
			return &command.PushCommand{
				Ui: ui,
			}, nil
		},
		"providers": func() (cli.Command, error) {
			return &command.ProvidersCommand{
				Ui: ui,
			}, nil
		},
		"repositories": func() (cli.Command, error) {
			return &command.RepositoriesCommand{
				Ui: ui,
			}, nil
		},
	}
}
