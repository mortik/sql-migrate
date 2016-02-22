package main

import (
	"flag"
	"strings"
)

type CreateCommand struct {
}

func (c *CreateCommand) Help() string {
	helpText := `
Usage: sql-migrate create [options] ...

  Creates a new Database Migration template.

Options:
  -config=dbconfig.yml   Configuration file to use.
  -env="development"     Environment.

`
	return strings.TrimSpace(helpText)
}

func (c *CreateCommand) Synopsis() string {
	return "Creates a new Database Migration template."
}

func (c *CreateCommand) Run(args []string) int {
	name := args[0]

	cmdFlags := flag.NewFlagSet("create", flag.ContinueOnError)
	cmdFlags.Usage = func() { ui.Output(c.Help()) }
	ConfigFlags(cmdFlags)

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := GenerateMigrationTemplate(name)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	return 0
}
