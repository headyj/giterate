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
		Version:  "0.2",
	}

	status, err := cli.Run()
	if err != nil {
		log.Fatalf("Cannot run command: %s", err)
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
	}
}
