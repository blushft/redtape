package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/blushft/redtape"
	"github.com/urfave/cli/v2"
)

func rolesCmd() *cli.Command {
	return &cli.Command{
		Name:        "roles",
		Usage:       "roles subcommand",
		Category:    "roles",
		Subcommands: []*cli.Command{roleBuildCmd()},
	}
}

func roleBuildCmd() *cli.Command {
	return &cli.Command{
		Name:     "build",
		Usage:    "build a new role",
		Category: "roles",
		Action:   roleBuildAction,
	}
}

func roleBuildAction(ctx *cli.Context) error {
	r, err := newRoleSurvey()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, string(b))

	return nil
}

type surveyRole struct {
	ID          string
	Name        string
	Description string
}

func newRoleSurvey() (*redtape.Role, error) {
	q := []*survey.Question{
		{
			Name:      "id",
			Prompt:    &survey.Input{Message: "role identifier:"},
			Validate:  survey.Required,
			Transform: survey.ToLower,
		},
		{
			Name:   "name",
			Prompt: &survey.Input{Message: "role friendly name:"},
		},
		{
			Name:   "description",
			Prompt: &survey.Input{Message: "role description"},
		},
	}

	a := surveyRole{}

	if err := survey.Ask(q, &a); err != nil {
		return nil, err
	}

	return &redtape.Role{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
	}, nil
}
