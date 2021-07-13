package main

import (
	"github.com/AlecAivazis/survey/v2"
)

func rangePrompt(prompt string) ([]string, error) {
	var a []string

	p := &survey.Input{
		Message: prompt,
		Help:    "x to continue",
	}

	for {
		var t string
		if err := survey.AskOne(p, &t); err != nil {
			return nil, err
		}

		if len(t) == 0 || t == "x" {
			break
		}

		a = append(a, t)
	}

	return a, nil
}
