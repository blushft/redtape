package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/blushft/redtape"
	"github.com/urfave/cli/v2"
)

func policyCmd() *cli.Command {
	return &cli.Command{
		Name:        "policy",
		Usage:       "policy subcommand",
		Category:    "policy",
		Subcommands: []*cli.Command{policyBuildCmd()},
	}
}

func policyBuildCmd() *cli.Command {
	return &cli.Command{
		Name:     "build",
		Usage:    "build a policy",
		Category: "policy",
		Action:   policyBuildAction,
	}
}

func policyBuildAction(ctx *cli.Context) error {
	p, err := newPolicySurvey()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, string(b))

	return nil
}

type surveyPolicy struct {
	Name        string
	Description string
	Roles       []string
	Resources   []string
	Actions     []string
	Scopes      []string
	Effect      string
}

func newPolicySurvey() (*redtape.PolicyOptions, error) {
	q := []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "policy name"},
			Validate:  survey.Required,
			Transform: survey.ToLower,
		},
		{
			Name:   "description",
			Prompt: &survey.Input{Message: "description"},
		},
	}

	a := surveyPolicy{}

	if err := survey.Ask(q, &a); err != nil {
		return nil, err
	}

	res, err := rangePrompt("match resources")
	if err != nil {
		return nil, err
	}

	a.Resources = res

	act, err := rangePrompt("match actions")
	if err != nil {
		return nil, err
	}

	a.Actions = act

	p := a.build()

	return &p, nil
}

func (p surveyPolicy) build() redtape.PolicyOptions {
	return redtape.NewPolicyOptions(
		redtape.PolicyName(p.Name),
		redtape.PolicyDescription(p.Description),
		redtape.SetResources(p.Resources...),
		redtape.SetActions(p.Actions...),
		redtape.SetScopes(p.Scopes...),
		redtape.SetPolicyEffect(p.Effect),
	)
}
